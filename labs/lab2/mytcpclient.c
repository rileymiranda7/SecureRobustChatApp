/* include libraries */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <unistd.h>
#include <netdb.h>

int main (int argc, char *argv[])
{
	if(argc!=3) {
		printf("Usage: %s <servername> <port>\n", argv[0]);
		exit(1);
	}
	printf("TCP Client program by Riley Miranda\n");
	char *servername = argv[1];
	char *port = argv[2];
	if ( strlen(servername) > 255 || strlen(port) > 5) {
		printf("Servername or port is too long.\
					Please try again!\n");
		exit(1);
	}
	printf("Servername= %s, port=%s\n", servername, port);
	int sockfd = socket(AF_INET, SOCK_STREAM, 0);
	if (sockfd < 0) {
		perror("ERROR opening socket");
		exit(2);
	}
	printf("A socket is opened!\n");
	//main code goes here
	struct addrinfo hints, *serveraddr;
	memset(&hints, 0, sizeof hints);
	hints.ai_family = AF_INET;
	hints.ai_socktype = SOCK_STREAM;
	int addr_lookup =
		getaddrinfo(servername, port, &hints, &serveraddr);
	if (addr_lookup != 0) {
		fprintf(stderr, "getaddrinfo: %s\n",
							gai_strerror(addr_lookup));
		exit(3);
	}
	int connected = connect(sockfd, serveraddr->ai_addr,
									serveraddr->ai_addrlen);
	if(connected < 0){
		perror("Cannot connect to the server\n");
		exit(4);
	}
	printf("Connected to the server %s at port %s\n",
					servername, port);
	freeaddrinfo(serveraddr);

	// Send data from user input
	int BUFFERSIZE = 1024; //define the size of the buffer
	char buffer[BUFFERSIZE]; //define the buffer
	bzero(buffer,BUFFERSIZE); //empty the buffer
	printf("Enter your message to send: ");
	fgets(buffer, BUFFERSIZE, stdin); // secure version of gets
	int byte_sent = send(sockfd,buffer, strlen(buffer), 0);

	// Receive data
	bzero(buffer,BUFFERSIZE);
	int byte_received = recv(sockfd, buffer, BUFFERSIZE, 0);
	if(byte_received < 0) {
		perror("Error in reading");
		exit(4);
	}
	printf("Received from server: %s", buffer);


	//put the below code at the end;
	close(sockfd);
	return 0;
}//end main function