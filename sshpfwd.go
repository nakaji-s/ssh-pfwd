package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"io"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHPortForward struct {
	SshAddr    string
	LocalAddr  string
	RemoteAddr string
}

// Get default location of a private key
func privateKeyPath() string {
	return os.Getenv("HOME") + "/.ssh/id_rsa"
}

// Get private key for ssh authentication
func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, _ := ioutil.ReadFile(keyPath)
	return ssh.ParsePrivateKey(buff)
}

// Get ssh client config for our connection
// SSH config will use 2 authentication strategies: by key and by password
func makeSshConfig(user, password string) (*ssh.ClientConfig, error) {
	key, err := parsePrivateKey(privateKeyPath())
	if err != nil {
		return nil, err
	}

	config := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &config, nil
}

// Handle local client connections and tunnel data to the remote serverq
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println("error while copy remote->local:", err)
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Println(err)
		}
		chDone <- true
	}()

	<-chDone
}

func (s SSHPortForward) Handle() {
	// Build SSH client configuration
	cfg, err := makeSshConfig(os.Getenv("USER"), "password")
	if err != nil {
		panic(err)
	}

	// Establish connection with SSH server
	conn, err := ssh.Dial("tcp", s.SshAddr, cfg)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Start local server to forward traffic to remote connection
	local, err := net.Listen("tcp", s.LocalAddr)
	if err != nil {
		panic(err)
	}
	defer local.Close()

	// Handle incoming connections
	for {
		client, err := local.Accept()
		if err != nil {
			panic(err)
		}

		// Establish connection with remote server
		remote, err := conn.Dial("tcp", s.RemoteAddr)
		if err != nil {
			panic(err)
		}

		handleClient(client, remote)
	}
}
