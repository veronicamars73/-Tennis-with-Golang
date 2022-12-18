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

func receiver(c chan int) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	bola_server := <-c
	println(bola_server)
	if bola_server == 0 {
		println("Ending point.")
		close(c)
	} else {
		receiver_ball := r.Intn(2)
		println("Receiver", receiver_ball)
		if receiver_ball == 0 {
			println("Ace point!")
			close(c)
		} else {
			c <- receiver_ball
		}
	}
}

func first(c chan int) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	bola_receiver := <-c
	println(bola_receiver)
	if bola_receiver == 0 {
		println("Ending point.")
		close(c)
	} else {
		server_ball := r.Intn(2)
		println("Server", server_ball)
		if server_ball == 0 {
			println("Server missed")
			c <- server_ball
		} else {
			println("Server returns ball. Point continues")
			c <- server_ball
		}
	}
}

func second(c chan int) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	bola_server := <-c
	println(bola_server)
	if bola_server == 0 {
		println("Ending point.")
		close(c)
	} else {
		receiver_ball := r.Intn(2)
		println("Receiver", receiver_ball)
		if receiver_ball == 0 {
			println("Receiver missed")
			c <- receiver_ball
		} else {
			println("Receiver returns ball. Point continues")
			c <- receiver_ball
		}
	}
}

func main() {

	quadra := make(chan int)

	go server(quadra)
	time.Sleep(1 * time.Second)
	go receiver(quadra)
	time.Sleep(1 * time.Second)
	_, status := <-quadra

	if status {
		fmt.Println("Point continues!")
		for status {
			go first(quadra)
			time.Sleep(1 * time.Second)
			go second(quadra)
			time.Sleep(1 * time.Second)
			_, status = <-quadra
		}
	} else {
		fmt.Println("Point ends!")
	}
}
