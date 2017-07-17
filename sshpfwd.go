package main

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHPortForward struct {
	SshAddr    string
	LocalAddr  string
	RemoteAddr string
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

	go func() {
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
	}()
}
