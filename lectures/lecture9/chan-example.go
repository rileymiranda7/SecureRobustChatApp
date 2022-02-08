package main

import "fmt"

func main(){
	i := make(chan int) //make channel i
	go func(){
		i <- 1 //write 1 to the channel
	}()
	//value_i := <- i //read from the channel
	fmt.Println(<-i)
}
