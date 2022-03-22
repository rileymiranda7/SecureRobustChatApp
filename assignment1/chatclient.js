var net = require('net');
var readlineSync = require('readline-sync');
const { getSystemErrorMap } = require('util');
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

	client.on("login failed", function(message) {
		loginsync();
	});

	client.on("data", data => {
		if (String.fromCharCode(...data) === "Login failed") {
			console.log("Login Failed: please retry...")
			loginsync();
		} else {
		process.stdout.write(data + "\n");
		}
	});

	client.on("error", function(err){
		console.log("Error");
		process.exit(2);
	});
	client.on("close", function(data){
		console.log("Connection has been disconnected");
		process.exit(3);
	});

	const keyboard = require('readline').createInterface({
		input: process.stdin,
		output: process.stout
	});
	keyboard.on('line', (input) => {
		console.log(`You typed: ${input}`);
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
	// Wait for user's response
	username = readlineSync.question('Username:');
	if (!inputValidated(username)) {
		console.log("Username must have at least 5 characters. Please try again!");
		loginsync();
		return;
	}
	//client.write(username);
	// Handle password text
	password = readlineSync.question('Password:', {
		hideEchoBack: true
	});
	if (!inputValidated(password)) {
		console.log("Password must have at least 5 characters. Please try again!");
		loginsync()
		return;
	}
	//client.write(password)
	var login = {Username: username, Password: password}
	client.write(JSON.stringify(login))
}

function inputValidated(loginCredential) {
	return (loginCredential.length > 4)
}