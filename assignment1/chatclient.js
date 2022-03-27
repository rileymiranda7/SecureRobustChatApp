var net = require('net');
var readlineSync = require('readline-sync');
//const { getSystemErrorMap } = require('util');
var username, password
if(process.argv.length != 4){
	console.log("Usage: node %s <host> <port>", process.argv[1]);
	process.exit(1);	
}

var host=process.argv[2];
var port=process.argv[3];

if(host.length >253 || port.length >5 ){
	console.log("Invalid host or port. Try again!\nUsage: node %s <port>", process.argv[1]);
	process.exit(1);	
}

var client = new net.Socket();
console.log("ChatClient.js developed by Riley Miranda, SecAD");
console.log("Connecting to: %s:%s", host, port);

client.connect(port,host, connected);

function connected(){
	console.log("Connected to: %s:%s", client.remoteAddress, client.remotePort);

	loginsync()

	/*client.on("login failed", function(message) {
		loginsync();
	});*/

	client.on("data", data => {
		var result = "";
		for(var i = 0; i < data.length; i++) {
			result += String.fromCharCode(parseInt(data[i])); // can't use spread operator
		}
		if (result === "Login failed") {
			console.log("Login Failed: please retry...")
			loginsync();
		} else {
		process.stdout.write(data + "\n");
		}
	});

	/*client.on("error", function(err){
		console.log("Error");
		process.exit(2);
	});*/
	client.on("close", function(data){
		console.log("Connection has been disconnected");
		process.exit(3);
	});

	const keyboard = require('readline').createInterface({
		input: process.stdin,
		output: process.stout
	});
	keyboard.on('line', (input) => {
		console.log(`You typed: ${input}\n`);
		//Some code here to handle input
		if(input === '.exit'){
			client.destroy();
			console.log('disconnected!');
			process.exit();
		} else {
		client.write(input);
		}
	});
}

function loginsync() {
	// Bug: readline-sync module is blocking. Multiple users can't log in at same time.
	// First user to connect has to finish logging in before the other's login data sends.
	// Possible solution: set timeout after 2 min if not logged in somehow?
	username = readlineSync.question('Username:');
	if (!inputValidated(username)) {
		loginsync();
		return;
	}
	//client.write(username);
	// Handle password text
	password = readlineSync.question('Password:', {
		hideEchoBack: true
	});
	if (!inputValidated(password)) {
		loginsync()
		return;
	}
	//client.write(password)
	var login = {Username: username, Password: password}
	client.write(JSON.stringify(login))
}

function inputValidated(loginCredential) {
	if (loginCredential.length > 4) {
		if (loginCredential.length < 1024) {
			return true
		} else {
			console.log("Input too long! Please try again:\n")
			return false;
		}
	} else {
		console.log("Input too short! Please try again:\n")
		return false
	}
	return false
}