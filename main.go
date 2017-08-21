package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	// Connection settings
	//sshAddr := "remote_ip:22"
	//sshAddr := "localhost:22"
	//localAddr := "127.0.0.1:5000"
	//remoteAddr := "127.0.0.1:8000"

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	config := SqliteConfig{db}
	defer db.Close()
	db.AutoMigrate(&Rule{})

	//config := InMemoryConfig{}
	//rule := Rule{Enable: false, Id: uuid.NewV4().String(), SSHPortForward: SSHPortForward{SshAddr: sshAddr, LocalAddr: localAddr, RemoteAddr: remoteAddr}}
	//config.AddRule(rule)

	server := Server{&config}
	go func() {
		server.Start()
	}()

	select {} // block
}
