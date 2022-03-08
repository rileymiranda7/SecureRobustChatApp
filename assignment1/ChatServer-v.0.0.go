/* Simple EchoServer in GoLang by Phu Phung, customized
by Riley Miranda for SecAD*/
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	//"encoding/json"
)

const BUFFERSIZE int = 1024

var allClient_conns = make(map[net.Conn]string)
var newclient = make(chan net.Conn)
var lostclient = make(chan net.Conn)

//var clientAuthenticated = make(chan net.Conn)

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
	fmt.Println("EchoServer in GoLang developed by Phu Phung, SecAD, revised by Riley Miranda")
	fmt.Printf("EchoServer is listening on port '%s' ...\n", port)
	go func() {
		for {
			client_conn, _ := server.Accept()
			newclient <- client_conn
		}
	}()
	for {
		select {
		case client_conn := <-newclient:
			if login(client_conn) {

				allClient_conns[client_conn] = client_conn.RemoteAddr().String()
				welcomemessage := fmt.Sprintf("A new client '%s' connected!\n# of connected clients: %d\n",
					client_conn.RemoteAddr().String(), len(allClient_conns))
				fmt.Println(welcomemessage)
				//broadcasting
				sendtoAll([]byte(welcomemessage))
				go client_goroutine(client_conn)
			}
		case client_conn := <-lostclient:
			//handling for the event
			delete(allClient_conns, client_conn)
			byemessage := fmt.Sprintf("Client '%s' is DISCONNECTED!\n# of connected clients: %d\n",
				client_conn.RemoteAddr().String(), len(allClient_conns))
			sendtoAll([]byte(byemessage))
			//case  client_conn := <- clientAuthenticated
		}
	}
}
func client_goroutine(client_conn net.Conn) {
	var buffer [BUFFERSIZE]byte
	for {
		byte_received, read_err := client_conn.Read(buffer[0:])
		if read_err != nil {
			fmt.Println("Error in receiving...")
			lostclient <- client_conn
			return
		}

		// input validation
		fmt.Printf("Received data: %sdata size = %d\n",
			string(buffer[:]), byte_received)
		if byte_received >= 7 && string(buffer[0:5]) == "login" {
			success_message := fmt.Sprintf("You typed: login\nReceived data: login data\n")
			_, write_err := client_conn.Write([]byte(success_message))
			if write_err == nil {
				fmt.Println("Sent data: login")
			}
		} else {
			err_message := fmt.Sprintf("You typed: %sReceived data: Non-login data\n", string(buffer[:]))
			_, write_err := client_conn.Write([]byte(err_message))
			if write_err == nil {
				fmt.Println("Sent data: Non-login data")
			}
		}
	}
}
func sendtoAll(data []byte) {
	for client_conn, _ := range allClient_conns {
		_, write_err := client_conn.Write(data)
		if write_err != nil {
			continue //move on next iteration
		}
	}
	fmt.Printf("Send data: %s to all clients!\n", data)
}
func sendto(data []byte, client_conn net.Conn) {
	_, write_err := client_conn.Write(data)
	if write_err != nil {
		fmt.Println("Error in receiving...")
		return
	}
}

// Reads input; returns false if input is empty (less than 3 bytes)
func readInput(client_conn net.Conn) (string, bool) {
	fmt.Println("@readInput()")
	var buffer [BUFFERSIZE]byte
	byte_received, read_err := client_conn.Read(buffer[0:])
	fmt.Printf("Read input of size %d\n", byte_received)
	if read_err != nil || byte_received < 3 {
		fmt.Println("Error in receiving or empty input...")
		return "error", true
	}
	return string(buffer[:]), false
}

func login(client_conn net.Conn) bool {
	fmt.Println("@login()")
	var username, password string
	var usernameInputEmpty bool = true
	var passwordInputEmpty bool = true
	return func() bool {
		sendto([]byte("   username: "), client_conn)
		username, usernameInputEmpty = readInput(client_conn)
		if !usernameInputEmpty { // if inputted username is not empty or there is error
			username = strings.Replace(username, " ", "", -1)
			fmt.Printf("username received: %s", username)
			sendto([]byte("   password: "), client_conn)
			password, passwordInputEmpty = readInput(client_conn)
			if !passwordInputEmpty {
				fmt.Printf("password received: %s", password)
				return checklogin(username, password, client_conn)
			} else {
				sendto([]byte("Error with entered password: Please try again\n"), client_conn)
				return login(client_conn)
			}
		} else { //username is empty
			sendto([]byte("Error with entered username: Please try again\n"), client_conn)
			return login(client_conn)
		}
	}()
	//return checklogin(username, password, client_conn)
}
func checklogin(username string, password string, client_conn net.Conn) bool {
	fmt.Println("@checklogin()")
	username = strings.Replace(username, " ", "", -1)
	password = strings.Replace(password, " ", "", -1)
	if username == "" || password == "" {
		fmt.Println("@checklogin() -> username or password empty")
		return login(client_conn)
	}
	var account Account
	account.username = username
	account.password = password
	return checkaccount(account, client_conn)
}
func checkaccount(account Account, client_conn net.Conn) bool {
	fmt.Println("@checkaccount()")
	if account.username != "" && account.password != "" && len(account.username) > 3 && len(account.password) > 3 {
		fmt.Println("@checkaccount() -> username and password received")
		return true
	} else {
		sendto([]byte("Invalid Username or Password...\n"), client_conn)
		return login(client_conn)
	}
}

// define Account struct
type Account struct {
	username string
	password string
}

//loginJson := `{"username": "bob", "password": "1234" }`
//var myLogin Login //var login is type Login
//json.Unmarshal([]byte(loginJson), &myLogin)
// Now can access username and password using dot notation:
// myLogin.username, myLogin.password
