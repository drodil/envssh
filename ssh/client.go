package ssh

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/drodil/envssh/util"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

// Client is used to run commands and set up environment
// variables in the remote.
type Client struct {
	sshClient    *ssh.Client
	envVariables map[string]string
	remote       *util.Remote
}

var logger = util.GetLogger()

// Connect connects to remote with given network, address and client configuration.
func Connect(network string, remoteAddr *util.Remote, config *ssh.ClientConfig) (*Client, error) {
	config.HostKeyCallback = checkHostKey

	logger.Println("Connecting to", remoteAddr.ToAddress())
	client, err := ssh.Dial(network, remoteAddr.ToAddress(), config)
	if err != nil {
		return nil, err
	}

	return &Client{
		sshClient:    client,
		envVariables: make(map[string]string),
		remote:       remoteAddr,
	}, nil
}

// TODO: Add support for private key authentication

// ConnectAuto tries to connect with all available methods.
// The sequence is SSH Agent, public key, password
func ConnectAuto(remote *util.Remote) (*Client, error) {
	authMethods := []ssh.AuthMethod{}
	sshAgent := getSshAgent()
	if sshAgent != nil {
		authMethods = append(authMethods, sshAgent)
	}

	config := &ssh.ClientConfig{
		User: remote.Username,
		Auth: authMethods,
	}

	client, err := Connect("tcp", remote, config)
	if err == nil {
		return client, nil
	}
	return ConnectWithPassword(remote)
}

func getSshAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

// ConnectWithSshAgent tries to connect to remote using sshagent. Only
// available on unix systems.
func ConnectWithSshAgent(remote *util.Remote) (*Client, error) {
	authMethods := []ssh.AuthMethod{}
	sshAgent := getSshAgent()
	if sshAgent != nil {
		authMethods = append(authMethods, sshAgent)
	}

	config := &ssh.ClientConfig{
		User: remote.Username,
		Auth: authMethods,
	}

	return Connect("tcp", remote, config)
}

func promptPassword(remote *util.Remote) ssh.AuthMethod {
	question := fmt.Sprint(remote.Username, "@", remote.Hostname, "'s password:")
	password := util.PromptPassword(question)
	return ssh.Password(password)
}

// ConnectWithPassword connect to the remote and prompts user for password.
func ConnectWithPassword(remote *util.Remote) (*Client, error) {
	var err error
	var client *Client

	for i := 0; i < 3; i++ {
		config := &ssh.ClientConfig{
			User: remote.Username,
			Auth: []ssh.AuthMethod{
				promptPassword(remote),
			},
		}
		client, err = Connect("tcp", remote, config)
		if err == nil {
			return client, nil
		}
		fmt.Println("Permission denied, please try again.")
	}
	return nil, err
}

// Disconnect disconnects the client.
func (client *Client) Disconnect() error {
	fmt.Println("Connection to", client.remote.Hostname, "closed.")
	return client.sshClient.Close()
}

// RunCommand runs single command in the remote.
func (client *Client) RunCommand(cmd string, stdIn, stdOut, stdErr *os.File) error {
	// TODO: Add support for STDOUT/STDERR
	session, err := client.sshClient.NewSession()
	if err != nil {
		return err
	}

	if stdIn != nil {
		session.Stdin = stdIn
	}
	if stdOut != nil {
		session.Stdout = stdOut
	}
	if stdErr != nil {
		session.Stderr = stdErr
	}

	logger.Println("Running command on remote ", cmd)
	if err := session.Run(cmd); err != nil {
		logger.Fatal("Failed to run command", err)
		return err
	}

	return nil
}

// MoveFileAtRemote moves file in remote from location to another.
func (client *Client) MoveFileAtRemote(from string, to string) error {
	cmd := fmt.Sprint("mv ", from, " ", to)
	return client.RunCommand(cmd, nil, nil, nil)
}

// CopyFileToRemote copies local file to remote over SSH.
func (client *Client) CopyFileToRemote(localFile string, remoteFile string) error {
	// TODO: Find a better way to do this. But not with SCP command.
	f, err := os.Open(localFile)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(content)

	logger.Println("Copying file", localFile, "to remote", remoteFile)
	cmd := fmt.Sprint("echo '", encoded, "' | base64 --decode > ", remoteFile)
	return client.RunCommand(cmd, nil, nil, nil)
}

// CopyFileFromRemote copies file from the remote to local over SSH.
func (client *Client) CopyFileFromRemote(remoteFile string, localFile string) error {
	// TODO: Find a better way to do this. But not with SCP command.
	file, err := ioutil.TempFile("", "envsshFileTransfer")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	cmd := fmt.Sprint("base64 ", remoteFile)
	err = client.RunCommand(cmd, nil, file, nil)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	decoded, err := base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return err
	}
	return ioutil.WriteFile(localFile, decoded, 0644)
}

// SetRemoteEnv sets remote environment variable that will be set when
// interactive session is started with StartInteractiveSession.
func (client *Client) SetRemoteEnv(name string, value string) {
	client.envVariables[name] = value
}

// SetRemoteEnvMap sets remote environment variables from map that will be set when
// interactive session is started with StartInteractiveSession.
func (client *Client) SetRemoteEnvMap(envVariables map[string]string) {
	for name, value := range envVariables {
		client.envVariables[name] = value
	}
}

// StartInteractiveSession starts interactive shell session with the remote.
func (client *Client) StartInteractiveSession() error {
	session, err := client.sshClient.NewSession()
	if err != nil {
		return err
	}

	// TODO: This only works for env variables that are listed in
	// sshd_config AcceptEnv. Maybe if pushing these as export after
	// connecting could allow setting other env variables as well
	for name, value := range client.envVariables {
		session.Setenv(name, value)
	}

	// TODO: Check modes
	modes := ssh.TerminalModes{
		ssh.ECHO: 1,
	}

	fd := int(Fd)
	state, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, state)

	w, h, err := term.GetSize(fd)
	if err != nil {
		// Default to 80x24
		w = 80
		h = 24
	}

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	if err := session.RequestPty(termType, h, w, modes); err != nil {
		// TODO: Fallback another PTY?
		return err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return err
	}

	// Handle terminal resize
	ch := make(chan os.Signal)
	signal.Notify(ch, ResizeEvent)
	go func() {
		for {
			s := <-ch
			switch s {
			case ResizeEvent:
				w, h, err = term.GetSize(fd)
				if err == nil {
					session.WindowChange(h, w)
				}
			}
		}
	}()

	// Wait until connection closes.
	if err := session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}
		return err
	}

	return nil
}

// TODO: Check if this could be replaced with https://pkg.go.dev/golang.org/x/crypto/ssh/knownhosts
func getKnownHosts(flag int, perm os.FileMode) (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath.Join(homeDir, ".ssh", "known_hosts"), flag, perm)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// TODO: Check if this could be replaced with https://pkg.go.dev/golang.org/x/crypto/ssh/knownhosts
func getHostKey(address string) ssh.PublicKey {
	file, err := getKnownHosts(os.O_RDONLY, 0744)

	if err != nil {
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], address) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				logger.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	return hostKey
}

// TODO: Check if this could be replaced with https://pkg.go.dev/golang.org/x/crypto/ssh/knownhosts
func addHostKey(address string, key ssh.PublicKey) error {
	file, err := getKnownHosts(os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0744)

	if err != nil {
		return err
	}

	defer file.Close()

	encoded := base64.StdEncoding.EncodeToString(key.Marshal())
	entry := fmt.Sprintln(address, key.Type(), encoded)
	if _, err := file.WriteString(entry); err != nil {
		return err
	}

	return nil
}

// TODO: Check if this could be replaced with https://pkg.go.dev/golang.org/x/crypto/ssh/knownhosts
func checkHostKey(hostname string, remote net.Addr, key ssh.PublicKey) error {
	hostnameWithoutPort := strings.Split(hostname, ":")[0]
	remoteWithoutPort := strings.Split(remote.String(), ":")[0]
	hostKey := getHostKey(hostnameWithoutPort)
	if hostKey == nil {
		fingerprint := ssh.FingerprintSHA256(key)
		fmt.Print("The authenticity of host '", hostnameWithoutPort, " (", remoteWithoutPort, ")' can't be established.\n")
		fmt.Print("ECDSA key fingerprint is ", fingerprint, "\n")
		answer := util.PromptAllowedString("Are you sure you want to continue connecting (yes/no)?", []string{"yes", "no"}, "no")
		if strings.Compare("yes", answer) != 0 {
			return errors.New("host key verification failed")
		}
		return addHostKey(hostnameWithoutPort, key)
	}
	return nil
}
