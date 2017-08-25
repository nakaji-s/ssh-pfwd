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
	SshAddr    string `validate:"required,gt=0"`
	LocalAddr  string `validate:"required,gt=0"`
	RemoteAddr string `validate:"required,gt=0"`
	Connected  bool
	conn       *ssh.Client
	local      net.Listener
	done       chan bool
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

func (s *SSHPortForward) Stop() {
	s.Connected = false

	if s.local != nil {
		s.done <- true
		s.local.Close()
		s.conn.Close()
	}
}

func (s *SSHPortForward) Start() {
	// Build SSH client configuration
	cfg, err := makeSshConfig(os.Getenv("USER"), "password")
	if err != nil {
		panic(err)
	}

	// Establish connection with SSH server
	s.conn, err = ssh.Dial("tcp", s.SshAddr, cfg)
	if err != nil {
		panic(err)
	}

	// Start local server to forward traffic to remote connection
	s.local, err = net.Listen("tcp", s.LocalAddr)
	if err != nil {
		panic(err)
	}

	s.Connected = true
	s.done = make(chan bool, 1)

	go func() {
		// Handle incoming connections
		for {
			client, err := s.local.Accept()
			if err != nil {
				select {
				case <-s.done:
					return
				default:
					panic(err)
				}
			}

			// Establish connection with remote server
			remote, err := s.conn.Dial("tcp", s.RemoteAddr)
			if err != nil {
				panic(err)
			}

			handleClient(client, remote)
		}
	}()
}
