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

func receiver(c chan int, is_serve bool, winner chan int) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	bola_server := <-c
	if bola_server == 0 {
		println("Ending point.")
		winner <- 1
		close(c)
	} else {
		receiver_ball := r.Intn(2)
		println("Receiver", receiver_ball)
		if is_serve {
			if receiver_ball == 0 {
				println("Ace point!")
				winner <- 0
				close(c)
			} else {
				c <- receiver_ball
				winner <- -1
			}
		} else {
			if receiver_ball == 0 {
				println("Receiver missed. Point server.")
				winner <- 0
				close(c)
			} else {
				c <- receiver_ball
				winner <- -1
			}
		}
	}
}

func score_calc(points int) {
	println(points)
	if points == 0 {
		println("0")
	} else {
		if points == 1 {
			println("15")
		} else {
			if points == 2 {
				println("30")
			} else {
				println("40")
			}
		}
	}
}

func main() {

	pontuacao_server := 0
	pontuacao_receiver := 0
	set_winner := -1

	for set_winner == -1 {
		quadra := make(chan int)
		point_winner := make(chan int)
		status := true
		serve := true
		var winner int
		for status {
			go server(quadra, serve)
			go receiver(quadra, serve, point_winner)
			time.Sleep(2 * time.Second)
			go func() {
				winner = <-point_winner
			}()
			_, status = <-quadra
			if winner != -1 && !status {
				if winner == 0 {
					println("The point winner was the server")
				} else {
					println("The point winner was the receiver")
				}
			}
			serve = false
		}
		if winner == 1 {
			pontuacao_receiver += 1
		} else {
			pontuacao_server += 1
		}
		if pontuacao_receiver == 3 && pontuacao_server == 3 {
			println("Deuce")
		} else {
			if pontuacao_receiver == 3 && pontuacao_server == 4 {
				println("Advantage server")
			} else {
				if pontuacao_receiver == 4 && pontuacao_server == 3 {
					println("Advantage receiver")
				} else {
					if pontuacao_receiver == 4 && pontuacao_server == 4 {
						println("Deuce")
						pontuacao_server = 3
						pontuacao_receiver = 3
					} else {
						if (pontuacao_receiver == 4 && pontuacao_server < 3) || (pontuacao_receiver == 5 && pontuacao_server == 3) {
							set_winner = 1
						} else {
							if (pontuacao_receiver < 3 && pontuacao_server == 4) || (pontuacao_receiver == 3 && pontuacao_server == 5) {
								set_winner = 0
							} else {
								println("Score:")
								print("Server score: ")
								score_calc(pontuacao_server)
								print("Receiver score: ")
								score_calc(pontuacao_receiver)
							}
						}
					}
				}
			}
		}
	}
	if set_winner == 0 {
		println("The set winner was the server")
	} else {
		println("The set winner was the receiver")
	}

}
