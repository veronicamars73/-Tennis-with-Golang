package main

import (
	"math/rand"
	"time"
)

func server(c chan int, is_serve bool) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	serve_result := r.Intn(2)
	if is_serve {
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
	} else {
		println("Server:", serve_result)
		if serve_result == 0 {
			println("Server missed.")
		}
	}

	c <- serve_result
}

func receiver(c chan int, is_serve bool) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	bola_server := <-c
	if bola_server == 0 {
		println("Ending point.")
		close(c)
	} else {
		receiver_ball := r.Intn(2)
		println("Receiver", receiver_ball)
		if is_serve {
			if receiver_ball == 0 {
				println("Ace point!")
				close(c)
			} else {
				c <- receiver_ball
			}
		} else {
			if receiver_ball == 0 {
				println("Receiver missed. Point server.")
				close(c)
			} else {
				c <- receiver_ball
			}
		}

	}
}

func main() {

	quadra := make(chan int)
	status := true
	serve := true

	for status {
		go server(quadra, serve)
		time.Sleep(1 * time.Second)
		go receiver(quadra, serve)
		time.Sleep(1 * time.Second)
		_, status = <-quadra
		serve = false
	}
}
