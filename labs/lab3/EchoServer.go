/* Simple EchoServer in GoLang by Phu Phung, customized 
by Riley Miranda for SecAD*/
package main

import (
	"fmt"
	"net"
	"os"
)

const BUFFERSIZE int = 1024
var allClient_conns = make(map[net.Conn]string)
var newclient = make(chan net.Conn)
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(0)
	}
	port := os.Args[1]
	if len(port) > 5 {
		fmt.Println("Invalid port value. Try again!")
		os.Exit(1)
	}
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Cannot listen on port '" + port + "'!\n")
		os.Exit(2)
	}
	fmt.Println("EchoServer in GoLang developed by Phu Phung, SecAD, revised by Your Name")
	fmt.Printf("EchoServer is listening on port '%s' ...\n", port)
	go func (){
		for {
			client_conn, _ := server.Accept()
			newclient <- client_conn
		}
	}()
	for {
		select{
				case client_conn := <- newclient:
					allClient_conns[client_conn]= client_conn.RemoteAddr().String()
					fmt.Printf("# of connected clients: %d\n", 
						len(allClient_conns))
					go client_goroutine(client_conn)
		}
	}
}
func client_goroutine(client_conn net.Conn){
	var buffer [BUFFERSIZE]byte
	for {
		byte_received, read_err := client_conn.Read(buffer[0:])
		if read_err != nil {
			fmt.Println("Error in receiving...")
			return
		}
		_, write_err := client_conn.Write(buffer[0:byte_received])
		if write_err != nil {
			fmt.Println("Error in sending...")
			return
		}
		fmt.Printf("Received data: %sEchoed back!\n", buffer)
	}
}
