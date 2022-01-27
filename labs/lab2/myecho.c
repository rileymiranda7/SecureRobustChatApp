/* include libraries */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main (int argc, char *argv[])
{
	if(strlen(argv[1]) > 126){
		printf("Input is too long");
		exit(1);
	}
	char buffer[126];
	strncpy(buffer,argv[1],126);
	printf("%s\n", argv[1]);
}