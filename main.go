package main

import (
	"fmt"
	"math/rand"
	"time"
)

func server(c chan int) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	serve_result := r.Intn(2)
	println("First serve:", serve_result)
	if serve_result == 0 {
		println("First fault!")
		rand.Seed(time.Now().UnixNano())
		serve_result = r.Intn(2)
		println("Second serve:", serve_result)
		if serve_result == 0 {
			println("Double Fault! Point to reciever.")
		}
	}
	c <- serve_result
}

func reciever(c chan int) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	bola_server := <-c
	println(bola_server)
	if bola_server == 0 {
		println("Ending point.")
		close(c)
	} else {
		reciever_ball := r.Intn(2)
		println("Recieve", reciever_ball)
		if reciever_ball == 0 {
			println("Ace point!")
			close(c)
		} else {
			println("Point continues")
		}
	}
}

func main() {

	quadra := make(chan int)

	go server(quadra)
	time.Sleep(1 * time.Second)
	go reciever(quadra)
	time.Sleep(1 * time.Second)
	_, status := <-quadra

	if status {
		fmt.Println("Point continues!")
		close(quadra)
	} else {
		fmt.Println("Point ends!")
	}
}
