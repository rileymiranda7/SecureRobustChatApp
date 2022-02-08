package main

import "fmt"

func main(){
	var i int
	func(){
		i = 1
	}()
	fmt.Println(i)
}
