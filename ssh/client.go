package ssh

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/drodil/envssh/util"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	sshClient *ssh.Client
}

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
		sshClient: client,
	}, nil
}

func (client *Client) Disconnect() error {
	return client.sshClient.Close()
}

// TODO: Add support for private key authentication

func ConnectWithPassword(address string, username string, password string) (*Client, error) {
	// TODO: Use RetryableAuthMethod and KeyboardInteractiveChallenge instead password parameter
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	return Connect("tcp", address, config)
}

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
