/* Simple EchoServer in GoLang by Phu Phung, customized
by Riley Miranda for SecAD*/
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	//"encoding/json"
)

const BUFFERSIZE int = 1024

// define Account struct
type Account struct {
	Username string
	Password string
}

// define User struct
type User struct {
	Username   string
	isLoggedIn bool
	Key        string
}

//loginJson := `{"username": "bob", "password": "1234" }`
//var myLogin Login //var login is type Login
//json.Unmarshal([]byte(loginJson), &myLogin)
// Now can access username and password using dot notation:
// myLogin.username, myLogin.password

var allLoggedIn_conns = make(map[net.Conn]interface{})
var newclient = make(chan net.Conn)
var lostclient = make(chan net.Conn)
var loggedInClient = make(chan net.Conn)

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
	fmt.Println("ChatServer in GoLang developed by Riley Miranda")
	fmt.Printf("ChatServer is listening on port '%s' ...\n", port)
	go func() {
		for {
			client_conn, _ := server.Accept()
			newclient <- client_conn
		}
	}()
	for {
		select {
		case client_conn := <-newclient:

			loginSuccessful, username := login(client_conn)
			if loginSuccessful {
				var justLoggedInUser = User{Username: username, isLoggedIn: true}
				allLoggedIn_conns[client_conn] = justLoggedInUser
				//allLoggedIn_conns[client_conn] = client_conn.RemoteAddr().String()
				welcomemessage := fmt.Sprintf("A new client '%s' connected!\n# of connected clients: %d\n",
					justLoggedInUser.Username, len(allLoggedIn_conns))
				fmt.Println(welcomemessage)
				//broadcasting
				sendtoAll([]byte(welcomemessage))
				go client_goroutine(client_conn)
			}

		case client_conn := <-lostclient:
			//handling for the event
			delete(allLoggedIn_conns, client_conn)
			byemessage := fmt.Sprintf("Client '%s' is DISCONNECTED!\n# of connected clients: %d\n",
				client_conn.RemoteAddr().String(), len(allLoggedIn_conns))
			sendtoAll([]byte(byemessage))
			//case  client_conn := <- clientAuthenticated
		}
	}
}
func client_goroutine(client_conn net.Conn) {
	//var buffer [BUFFERSIZE]byte
	var username = allLoggedIn_conns[client_conn].(User).Username
	for {
		menu := fmt.Sprintf("Type the number of the operation you would like to perform:\n" +
			"1) Get List of Online Users [1 + Enter/Return]\n2)  Send message to all online users " +
			"[2 + Enter]\n3)  Send private message [3 + Enter]--Type 'help' to display options again\n")
		sendto([]byte(menu), client_conn)
		sendto([]byte("  Option: "), client_conn)
		var optionNum = readInput(client_conn)
		fmt.Printf("optionNum: %s\n", optionNum)
		switch optionNum {
		case "1":
			sendUserList(client_conn)
		case "2":
			sendto([]byte("Type message:"), client_conn)
			sendtoAll([]byte("[" + username + "]: " + readInput(client_conn) + "\n"))
		case "3":
			sendUserList(client_conn)
			sendto([]byte("Type username of receiver:"), client_conn)
			receiver := readInput(client_conn)
			sendto([]byte("Type message to send to "+receiver), client_conn)
			message := readInput(client_conn)
			sendPrivateM(receiver, username, message)
		case "help":
			continue
		default:
			//sendto([]byte("Invalid Option"), client_conn)
			sendto([]byte("Invalid Option\n"), client_conn)
			continue
		}
		/*byte_received, read_err := client_conn.Read(buffer[0:])
		if read_err != nil {
			fmt.Println("Error in receiving...")
			lostclient <- client_conn
			return
		}
		// input validation
		fmt.Printf("Received data: %sdata size = %d\n",
			string(buffer[:]), byte_received)*/
	}
}
func sendtoAll(data []byte) {
	for client_conn, _ := range allLoggedIn_conns {
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
func sendUserList(client_conn net.Conn) {
	var userlist = "Online users: "
	for client_conn, _ := range allLoggedIn_conns {
		user := allLoggedIn_conns[client_conn].(User)
		userlist += user.Username + ", "
	}
	sendto([]byte(userlist+"\n"), client_conn)
}
func sendPrivateM(receiver string, sender string, message string) {
	for client_conn, _ := range allLoggedIn_conns {
		if allLoggedIn_conns[client_conn].(User).Username == receiver {
			sendto([]byte("[private message from "+sender+"]: "+message+"\n"), client_conn)
		}
	}
}

// Reads input; returns false if input is empty (less than 3 bytes)
func readInput(client_conn net.Conn) string {
	//fmt.Println("@readInput()")
	var buffer [BUFFERSIZE]byte
	byte_received, read_err := client_conn.Read(buffer[0:])
	//fmt.Printf("Read input of size %d\n", byte_received)
	if read_err != nil /*|| byte_received < 3*/ {
		fmt.Println("Error in receiving or empty input...")
		return "error"
	}
	return string(buffer[:byte_received])
}

/**
Get login data and call checklogin()
If checklogin() returns true send client to authenticated channel to put on list
Else send error message back to client and call login() again
*/
func login(client_conn net.Conn) (bool, string) {
	//fmt.Println("@login()")
	var loginJSONstring = readInput(client_conn)
	fmt.Printf("Received login JSON: %s\n", loginJSONstring)
	loginSuccessful, message := checklogin(loginJSONstring, client_conn)
	if !loginSuccessful {
		return login(client_conn)
	} else {
		return loginSuccessful, message
	}
}

/**
Parse the username and password
Call checkaccount() to see if valid
*/
func checklogin(loginJSONstring string, client_conn net.Conn) (bool, string) {
	//fmt.Println("@checklogin()")
	var account Account
	err := json.Unmarshal([]byte(loginJSONstring), &account)
	if err != nil /*|| account.username == "" || account.password == ""*/ {
		fmt.Printf("JSON parsing error: %s\n", err)
		return false, "error"
	}
	//fmt.Printf("account.username: %s\naccount.password: %s\n", account.Username, account.Password)
	return checkaccount(account, client_conn)
}

/**
Compare username/password with database
*/
func checkaccount(account Account, client_conn net.Conn) (bool, string) {
	//fmt.Println("@checkaccount()")
	if (account.Username == "riley" && account.Password == "12345") ||
		(account.Username == "user0" && account.Password == "67890") ||
		(account.Username == "user00" && account.Password == "123456789") {
		//fmt.Println("@checkaccount() -> username and password received")
		return true, account.Username
	} else {
		//sendto([]byte("Invalid Username or Password...\n"), client_conn)
		fmt.Printf("Login <%s, %s> failed.\n", account.Username, account.Password)
		sendto([]byte("Login failed"), client_conn)
		return login(client_conn)
	}
}
