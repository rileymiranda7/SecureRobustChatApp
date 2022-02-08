package main

import "fmt"

func main(){
	var i int
	go func(){
		i = 1
	}()
	fmt.Println(i)
}
