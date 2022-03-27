/* ChatServer
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

var allLoggedIn_conns = make(map[net.Conn]interface{})
var newclient = make(chan net.Conn)
var lostclient = make(chan net.Conn)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(0)
	}
	port := os.Args[1]
	if len(port) > 5 {
		fmt.Println(">>Invalid port value. Try again!")
		os.Exit(1)
	}
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf(">>Cannot listen on port '" + port + "'!\n")
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
			fmt.Println("Debug>@case newclient")
			loginSuccessful, username := login(client_conn)
			if loginSuccessful {
				var justLoggedInUser = User{Username: username, isLoggedIn: true}
				allLoggedIn_conns[client_conn] = justLoggedInUser
				welcomemessage := fmt.Sprintf("A new client '%s' connected!\n# of connected clients: %d\n",
					justLoggedInUser.Username, len(allLoggedIn_conns))
				fmt.Println(welcomemessage)
				sendtoAll([]byte(welcomemessage))
				go client_goroutine(client_conn)
			}

		case client_conn := <-lostclient:
			fmt.Println("@lostclient")
			delete(allLoggedIn_conns, client_conn)
			byemessage := fmt.Sprintf("Client '%s' is DISCONNECTED!\n# of connected clients: %d\n",
				client_conn.RemoteAddr().String(), len(allLoggedIn_conns))
			var userlist = "Online users: "
			for client_conn, _ := range allLoggedIn_conns {
				user := allLoggedIn_conns[client_conn].(User)
				userlist += user.Username + ", "
			}
			fmt.Printf("%s\n", userlist)
			sendtoAll([]byte(byemessage))
		}
	}
}
func client_goroutine(client_conn net.Conn) {
	var username = allLoggedIn_conns[client_conn].(User).Username
	for {
		menu := fmt.Sprintf("Type the number of the operation you would like to perform:\n" +
			"1) Get List of Online Users [1 + Enter]\n2)  Send message to all online users " +
			"[2 + Enter]\n3)  Send private message [3 + Enter]\n.exit) Exit Chat Server [.exit + Enter]\n" +
			"--Type 'help' to display options again\n")
		sendto([]byte(menu), client_conn)
		sendto([]byte("  Option: "), client_conn)
		var optionNum, read_err = readInput(client_conn)
		if read_err != nil {
			break
		}
		fmt.Printf("optionNum: %s\n", optionNum)
		switch optionNum {
		case "1":
			sendUserList(client_conn)
		case "2":
			sendto([]byte("Type message to send to all:"), client_conn)
			var input, read_err = readInput(client_conn)
			if read_err != nil {
				break
			}
			sendtoAll([]byte("[" + username + "]: " + input + "\n"))
		case "3":
			sendUserList(client_conn)
			sendto([]byte("Type username of receiver:"), client_conn)
			receiver, read_err := readInput(client_conn)
			if read_err != nil {
				break
			}
			if !userIsOnline(client_conn, receiver) {
				continue
			} else {
				sendto([]byte("Type message to send to "+receiver), client_conn)
				message, read_err := readInput(client_conn)
				if read_err != nil {
					break
				}
				sendPrivateM(receiver, username, message)
			}
		case "help":
			continue
		default:
			sendto([]byte(">>Invalid Option\n"), client_conn)
			continue
		}
	}
}
func sendtoAll(data []byte) {
	fmt.Println("@sendtoAll")
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
		fmt.Println("@sendto(): Error in receiving...")
		return
	}
}
func sendUserList(client_conn net.Conn) {
	var userlist = "Online users: "
	var userlistSlice []string
	for client_conn, _ := range allLoggedIn_conns {
		user := allLoggedIn_conns[client_conn].(User)
		userlistSlice = append(userlistSlice, user.Username)
	}
	userlistSlice = removeDuplicateStrings(userlistSlice)
	for _, username := range userlistSlice {
		userlist += username + ", "
	}
	sendto([]byte(userlist+"\n\n"), client_conn)
}
func removeDuplicateStrings(slice []string) []string {
	strmap := make(map[string]bool) // initialize map where key will be values of str slice
	list := []string{}
	for _, item := range slice {
		if _, val := strmap[item]; !val { // initialize bool val from map that tests if item from slice is in map
			strmap[item] = true // if val not in map add to map and append to list; if is in map do nothing
			list = append(list, item)
		}
	}
	return list
}
func userIsOnline(client_conn net.Conn, user string) bool {
	var messageWasSent = false
	for client_conn, _ := range allLoggedIn_conns {
		if allLoggedIn_conns[client_conn].(User).Username == user {
			messageWasSent = true
		}
	}
	if messageWasSent {
		return true
	} else {
		sendto([]byte(">>Selected user is not online!\n\n"), client_conn)
		return false
	}
}
func sendPrivateM(receiver string, sender string, message string) {
	for client_conn, _ := range allLoggedIn_conns {
		if allLoggedIn_conns[client_conn].(User).Username == receiver {
			sendto([]byte("[private message from "+sender+"]: "+message+"\n\n"), client_conn)
		}
	}
}
func readInput(client_conn net.Conn) (string, error) {
	var buffer [BUFFERSIZE]byte
	byte_received, read_err := client_conn.Read(buffer[0:])
	if read_err != nil {
		fmt.Println("@readInput(): Error in receiving or empty input...")
		lostclient <- client_conn // client disconnects
		return "error", read_err
	}
	return string(buffer[:byte_received]), nil
}

/**
Get login data and call checklogin()
Returns true or calls login again
*/
func login(client_conn net.Conn) (bool, string) {
	var loginJSONstring, _ = readInput(client_conn)
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
	var account Account
	err := json.Unmarshal([]byte(loginJSONstring), &account)
	if err != nil {
		fmt.Printf("JSON parsing error: %s\n", err)
		return false, "error"
	}
	return checkaccount(account, client_conn)
}

/**
Check database and either return true or retry login
*/
func checkaccount(account Account, client_conn net.Conn) (bool, string) {
	if (account.Username == "riley" && account.Password == "12345") ||
		(account.Username == "user0" && account.Password == "67890") ||
		(account.Username == "user00" && account.Password == "123456789") {
		return true, account.Username
	} else {
		fmt.Printf("Login <%s, %s> failed.\n", account.Username, account.Password)
		sendto([]byte("Login failed"), client_conn)
		return login(client_conn)
	}
}
