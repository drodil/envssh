package ssh

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/drodil/envssh/util"
	"github.com/shiena/ansicolor"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	sshClient    *ssh.Client
	envVariables map[string]string
}

// Connects to remote with given network, address and client configuration.
func Connect(network string, address string, config *ssh.ClientConfig) (*Client, error) {
	config.HostKeyCallback = checkHostKey
	// TODO: Use some struct for host/port combination
	if !strings.Contains(address, ":") {
		address = address + ":22"
	}

	client, err := ssh.Dial(network, address, config)
	if err != nil {
		return nil, err
	}

	return &Client{
		sshClient:    client,
		envVariables: make(map[string]string),
	}, nil
}

// TODO: Add support for private key authentication
// TODO: Add support for SSHAgent
// TODO: Add support for auto connect with first available AuthMethod (sshagent, key, password)

// Connects to the remote with given username. Prompts user for password.
func ConnectWithPassword(address string, username string) (*Client, error) {
	// TODO: Add support to retry password input
	question := fmt.Sprint(username, "@", address, "'s password:")
	password := util.PromptPassword(question)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	return Connect("tcp", address, config)
}

// Disconnects the client.
func (client *Client) Disconnect() error {
	return client.sshClient.Close()
}

// Runs single command in the remote.
func (client *Client) RunCommand(cmd string) error {
	// TODO: Might require STD mapping? Test this plz.
	session, err := client.sshClient.NewSession()
	if err != nil {
		return err
	}

	// TODO: Use session.CombinedOutput for error code ?
	if err := session.Run(cmd); err != nil {
		return err
	}

	return nil
}

// Moves file in remote from location to another.
func (client *Client) MoveFileAtRemote(from string, to string) error {
	cmd := fmt.Sprint("mv ", from, " ", to)
	return client.RunCommand(cmd)
}

// Copies local file to remote over SSH.
func (client *Client) CopyFileToRemote(localFile string, remoteFile string) error {
	// TODO: Find a better way to do this. But not with SCP command.
	cmd := fmt.Sprint("\"cat > ", remoteFile, "\" < ", localFile)
	return client.RunCommand(cmd)
}

// Copies file from the remote to local over SSH.
func (client *Client) CopyFileFromRemote(remoteFile string, localFile string) error {
	// TODO: Find a better way to do this. But not with SCP command.
	cmd := fmt.Sprint("\"cat ", remoteFile, "\" > ", localFile)
	return client.RunCommand(cmd)
}

// Sets remote environment variable that will be set when
// interactive session is started with StartInteractiveSession.
func (client *Client) SetRemoteEnv(name string, value string) {
	client.envVariables[name] = value
}

// Starts interactive session with the remote.
func (client *Client) StartInteractiveSession() error {
	session, err := client.sshClient.NewSession()
	if err != nil {
		return err
	}

	for name, value := range client.envVariables {
		session.Setenv(name, value)
	}

	session.Stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
	session.Stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
	in, _ := session.StdinPipe()

	// TODO: Check modes
	modes := ssh.TerminalModes{
		ssh.ECHO:  0,
		ssh.IGNCR: 1,
	}

	// TODO: Get size of the Pty from current terminal
	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
		// TODO: Fallback another PTY?
		return err
	}

	if err := session.Shell(); err != nil {
		return err
	}

	// CTRL + C
	// TODO: Handler more signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c
			fmt.Println("^C")
			fmt.Fprint(in, "\n")
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')
		_, err := fmt.Fprint(in, str)
		// TODO: This continues correctly after server has disconnected session BUT
		// requires extra input from user.. Maybe use goroutine to check this?
		if err != nil {
			break
		}
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
		log.Fatal(err)
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
				log.Fatalf("error parsing %q: %v", fields[2], err)
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
		log.Fatal(err)
		return err
	}

	defer file.Close()

	encoded := base64.StdEncoding.EncodeToString(key.Marshal())
	entry := fmt.Sprintln(address, key.Type(), encoded)
	if _, err := file.WriteString(entry); err != nil {
		return nil
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
			return errors.New("Host key verification failed.")
		}
		return addHostKey(hostnameWithoutPort, key)
	}
	return nil
}
