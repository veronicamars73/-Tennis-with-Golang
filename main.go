package main

import (
	"math/rand"
	"time"
	"fmt"
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

func score_calc(points int) string {
	if points == 0 {
		return "00"
	} else {
		if points == 1 {
			return "15"
		} else {
			if points == 2 {
				return "30"
			} else {
				return "40"
			}
		}
	}
}

func show_score(pontuacao_server int, pontuacao_receiver int, games_server int, games_receiver int, sets_server int, sets_receiver int) {
	println("+-----------------------  SCORE  -------------------------+")
	println("|                                                         |")
	println("|               POINTS     -     GAMES     -     SETS     |")
	println("| Server:        ", score_calc(pontuacao_server), "             ", games_server, "            ", sets_server, "      |",  )
	println("|                                                         |")
	println("| Receiver:      ", score_calc(pontuacao_receiver), "             ", games_receiver, "            ", sets_receiver, "      |",  )
	println("|                                                         |")
	println("+---------------------------------------------------------+")
}

func main() {

	println(" ")
	println(" 00000000 0000000 000    00  00   000000")
	println("    00    00      00 0   00  00  00")
	println("    00    00000   00  0  00  00   00000")
	println("    00    00      00   0 00  00       00")
	println("    00    0000000 00    000  00  000000")
	println(" ")
	println("By: Alaide Lisandra e José Victor")
	println(" ")

	sets_server := 0
	sets_receiver := 0

	var numberOfSet int
	var numberOfGamesInSet int

	winner_match := -1

	print("Entry number of sets in the match: ")
	fmt.Scanln(&numberOfSet)

	print("Entry number of games in the set: ")
	fmt.Scanln(&numberOfGamesInSet)

	for winner_match == -1 {

		games_server := 0
		games_receiver := 0

		for winner_set := false; winner_set == false; winner_set = (games_server == numberOfGamesInSet || games_receiver == numberOfGamesInSet) {

		//Pontuação inicial do game do sacador
		pontuacao_server := 0
		//Pontuação inicial do game do receptor
		pontuacao_receiver := 0

		//Set com vencedor indefinido
		set_winner_game := -1
		
		//Laço de repetição para jogadas continuarem enquanto o vencedor do SET não for definido
		for set_winner_game == -1 {
				quadra := make(chan int)
				point_winner := make(chan int)
				status := true
				serve := true
				var winner_game int
				for status {
					go server(quadra, serve)
					go receiver(quadra, serve, point_winner)
					time.Sleep(2 * time.Second)
					go func() {
						winner_game = <-point_winner
					}()
					_, status = <-quadra
					if winner_game != -1 && !status {
						if winner_game == 0 {
							println("The point winner was the server")
						} else {
							println("The point winner was the receiver")
						}
					}
					serve = false
				}
				if winner_game == 1 {
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
									set_winner_game = 1
								} else {
									if (pontuacao_receiver < 3 && pontuacao_server == 4) || (pontuacao_receiver == 3 && pontuacao_server == 5) {
										set_winner_game = 0
									} else {
										show_score(pontuacao_server, pontuacao_receiver, games_server, games_receiver, sets_server, sets_receiver)
									}
								}
							}
						}
					}
				}
			}
			if set_winner_game == 0 {
				println("The game winner was the server")
				games_server++
			} else {
				println("The game winner was the receiver")
				games_receiver++
				
			}

			show_score(pontuacao_server, pontuacao_receiver, games_server, games_receiver, sets_server, sets_receiver)

		}
		if(games_server == numberOfGamesInSet){
					sets_server++
					println("The set winner was the server")
					if(sets_server >= numberOfSet && sets_server > sets_receiver ){
						winner_match = 0
					}
		 } else {
			if(games_receiver == numberOfGamesInSet){
				sets_receiver++
				println("The set winner was the receiver")
				if(sets_receiver >= numberOfSet && sets_server < sets_receiver){
					winner_match = 1
				}
			}
		}
	}

	if(winner_match == 0){
		println("The match winner was the server")
	} else {
		println("The match winner was the receiver")
	}

}
