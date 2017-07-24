package main

func main() {
	// Connection settings
	//sshAddr := "remote_ip:22"
	sshAddr := "localhost:22"
	localAddr := "127.0.0.1:5000"
	remoteAddr := "127.0.0.1:8000"

	config := Config{}
	rule := Rule{SshAddr: sshAddr, LocalAddr: localAddr, RemoteAddr: remoteAddr}
	config.AddRule(rule)

	server := Server{config}
	go func() {
		server.Start()
	}()

	pfwd := SSHPortForward{SshAddr: sshAddr, LocalAddr: localAddr, RemoteAddr: remoteAddr}
	go func() {
		pfwd.Handle()
	}()

	select {} // block
}
