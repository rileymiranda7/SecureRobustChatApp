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
	if ( strlen(servername) > 253 || strlen(port) > 5) {
		printf("Servername or port is too long.\
					Please try again!\n");
		exit(2);
	}
	printf("Servername= %s, port=%s\n", servername, port);
	int sockfd = socket(AF_INET, SOCK_STREAM, 0);
	if (sockfd < 0) {
		perror("ERROR opening socket");
		exit(sockfd);
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
		exit(3);
	}else
		printf("Connected to the server %s at port %s\n",
					servername, port);
	freeaddrinfo(serveraddr);


	//put the below code at the end;
	close(sockfd);
}
