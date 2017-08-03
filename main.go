package main

import "github.com/satori/go.uuid"

func main() {
	// Connection settings
	//sshAddr := "remote_ip:22"
	sshAddr := "localhost:22"
	localAddr := "127.0.0.1:5000"
	remoteAddr := "127.0.0.1:8000"

	config := InMemoryConfig{}
	rule := Rule{Enable: true, Id: uuid.NewV4().String(), SSHPortForward: SSHPortForward{SshAddr: sshAddr, LocalAddr: localAddr, RemoteAddr: remoteAddr}}
	config.AddRule(rule)

	server := Server{&config}
	go func() {
		server.Start()
	}()

	pfwd := rule.SSHPortForward
	go func() {
		pfwd.Handle()
	}()

	select {} // block
}
