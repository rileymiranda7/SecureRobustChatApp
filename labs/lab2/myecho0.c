/* include libraries */
#include <stdio.h>
#include <stdlib.h>

int main (int argc, char *argv[])
{
	//printf(argv[1]); Not secure not robust

	//printf("%s\n", argv[1]); // Secure but not robust (calling with no
	// arguments gives segmentation fault)

	if (argv[1]) { // more secure and more robust
		prinf("%s\n", argv[1]);
	} else {
		printf("Usage: %s <input>\n", argv[0]);
	}
}
