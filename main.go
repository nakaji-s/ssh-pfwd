package main

func main() {
	//config := Config{}

	server := Server{}
	go func() {
		server.Start()
	}()

	// Connection settings
	//sshAddr := "remote_ip:22"
	sshAddr := "localhost:22"
	localAddr := "127.0.0.1:5000"
	remoteAddr := "127.0.0.1:8000"

	pfwd := SSHPortForward{SshAddr: sshAddr, LocalAddr: localAddr, RemoteAddr: remoteAddr}
	go func() {
		pfwd.Handle()
	}()

	select {} // block
}
