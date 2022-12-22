package main

import (
	"math/rand"
	"time"
	"fmt"
)


/* 

--> Função Server:
		Recebe como parametro o canal que marca o resultado da jogada
		e um boleano para saber se saque ou não.

*/
func server(c chan int, is_serve bool) {
	/*
	  Sorteia um número aleatório para representar a jogada do sacador
		--> Número 0 representa o erro
	 	--> Número 1 representa o acerto
	*/
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	serve_result := r.Intn(2)

	//Verifica se é o sacador do jogo
	if is_serve {
		//Sacador executa sua primeira jogada
		println("First serve:", serve_result)
		//Verifica se o saque falhou
		if serve_result == 0 {
			/* 
				Caso tenha errado o saque, informa a primeira falta e sorteia 
				um novo número (0 ou 1) para simular nova jogada, verificando 
				novamente a ocorrencia de erros. Caso tenha errado, o recebedor
				marca o ponto.
			*/
			println("First fault!")
			rand.Seed(time.Now().UnixNano())
			serve_result = r.Intn(2)
			println("Second serve:", serve_result)
			if serve_result == 0 {
				println("Double Fault! Point to reciever.")
			}
		}
	} else {
		/* 
			Caso o sacador não seja o jogador da vez, é verificado
			 se ele conseguiu rebater a bola que chegou até ele.
		*/
		println("Server:", serve_result)
		if serve_result == 0 {
			println("Server missed.")
		}
	}

	//Escreve o resultado da jogada no canal de quadra
	c <- serve_result
}
/*
	--> Função Receiver:
		Recebe como parametro o canal que marca o resultado da jogada
		e um boleano para saber se o sacador acertou a bola.

*/
func receiver(c chan int, is_serve bool, winner chan int) {
	/*
	  Sorteia um número aleatório para representar a jogada do recebedor
		--> Número 0 representa o erro
	 	--> Número 1 representa o acerto
	*/
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	//Receber resultado da última jogada
	bola_server := <-c

	//Verifica se o sacador errou o ponto, então é ponto do receptor
	if bola_server == 0 {
		println("Ending point.")
		//Marca jogador como vencedor do ponto
		winner <- 1
		//Fecha o canal
		close(c)
	} else {
		/* 
			Caso o sacador não tenha errado o ponto, então, 
			sorteia o número para simular a rodada.
		*/
		receiver_ball := r.Intn(2)
		println("Receiver", receiver_ball)
		if is_serve {
			if receiver_ball == 0 {
				//Não conseguiu receber a bola e perdeu ponto
				println("Ace point!")
				//Marca o ponto e fecha o canal
				winner <- 0
				close(c)
			} else {
				//Escreve o resultado da jogada no canal
				c <- receiver_ball
				winner <- -1
			}
		} else {
			/* 
			Caso o sacador não seja o jogador da vez, é verificado
			 se ele conseguiu rebater a bola que chegou até ele.
			*/
			if receiver_ball == 0 {
				println("Receiver missed. Point server.")
				//Marca o ponto e fecha o canal
				winner <- 0
				close(c)
			} else {
				c <- receiver_ball
				winner <- -1
			}
		}
	}
}
/*
	Função Score: recebe o número de pontos na partida e retorna o valor representado.
*/
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


//Exibe o placar dois jogador
func show_score(points_server int, points_receiver int, games_server int, games_receiver int, sets_server int, sets_receiver int) {
	println("+-----------------------  SCORE  -------------------------+")
	println("|                                                         |")
	println("|               POINTS     -     GAMES     -     SETS     |")
	println("| Server:        ", score_calc(points_server), "             ", games_server, "            ", sets_server, "      |",  )
	println("|                                                         |")
	println("| Receiver:      ", score_calc(points_receiver), "             ", games_receiver, "            ", sets_receiver, "      |",  )
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

	//Numero de sets vencidos inicialmente pelos jogadores
	sets_server := 0
	sets_receiver := 0


	//Numero de sets do jogo
	var numberOfSet int
	//Numero de games por set de jogo
	var numberOfGamesInSet int

	//Jogo com vencedor indefinido (Match indefinido)
	winner_match := -1

	//Recebe pela entrada padrão o numero de sets no match
	print("Entry number of sets in the match: ")
	fmt.Scanln(&numberOfSet)

	//Recebe pela entrada padrão o numero de games por set
	print("Entry number of games in the set: ")
	fmt.Scanln(&numberOfGamesInSet)

	//Laço de repetição que executa enquanto não tiver vencedor do match
	for winner_match == -1 {

		//Quantidade de games vencidos por cada jogador (server e receiver)
		games_server := 0
		games_receiver := 0

		/* 
			Laço de repetição que executa enquanto um dos jogadores não
		 	atingir o numero de games necessarios para vender o SET. Ou seja,
			enquanto o vencedor do SET não for definido.
		*/
		for winner_set := false; winner_set == false; winner_set = (games_server == numberOfGamesInSet || games_receiver == numberOfGamesInSet) {

		//Pontuação inicial do game do sacador
		points_server := 0
		//Pontuação inicial do game do receptor
		points_receiver := 0

		//Game com vencedor indefinido
		set_winner_game := -1
		
		//Laço de repetição para jogadas continuarem enquanto o vencedor do game não for definido
		for set_winner_game == -1 {
				//Cria o canal pra marcar as jogadas de cada jogador 
				quadra := make(chan int)
				//Cria canal pra pegar o vencedor
				point_winner := make(chan int)
				/* 
					Variáveis boleanas para informar se o game ainda esta acontecendo (status)
					e para para informar se o jogador da vez foi o sacador (server) 
				*/
				status := true
				serve := true
				//Variável para armazenar o vencedor do game
				var winner_game int
				//Laço de repetição para simular as jogadas 
				for status {
					go server(quadra, serve)
					go receiver(quadra, serve, point_winner)
					time.Sleep(2 * time.Second)
					go func() {
						winner_game = <-point_winner
					}()
					_, status = <-quadra
					//Após a simulação, verifica se foi marcado ponto
					if winner_game != -1 && !status {
						if winner_game == 0 {
							println("The point winner was the server")
						} else {
							println("The point winner was the receiver")
						}
					}
					serve = false
				}
				//Atribui pontos a quem acertou a jogada
				if winner_game == 1 {
					points_receiver += 1
				} else {
					points_server += 1
				}
				/*
					No jogo de tenis, um game pode ser vencido ao marcar 4 pontos:

					1 ponto = 15 pontos
					2 pontos = 30 pontos
					3 pontos = 45 pontos
					4 pontos = GAME

					Por condição, a diferença de pontos entre os jogadores deve ser de 2 pontos.
					As linhas abaixo realizam essa verificação, não permitindo vencer com apenas
					um ponto de diferença.
				*/
				if points_receiver == 3 && points_server == 3 {
					println("Deuce")
				} else {
					if points_receiver == 3 && points_server == 4 {
						println("Advantage server")
					} else {
						if points_receiver == 4 && points_server == 3 {
							println("Advantage receiver")
						} else {
							if points_receiver == 4 && points_server == 4 {
								println("Deuce")
								points_server = 3
								points_receiver = 3
							} else {
								/*
									Aqui é virificado se a quantidade de ponto é suficiente para vencer
									e se a diferença é maior ou igual a 2. Caso seja, define o vencedor 
									e exibe o placar do jogo.

								*/
								if (points_receiver == 4 && points_server < 3) || (points_receiver == 5 && points_server == 3) {
									set_winner_game = 1
								} else {
									if (points_receiver < 3 && points_server == 4) || (points_receiver == 3 && points_server == 5) {
										set_winner_game = 0
									} else {
										show_score(points_server, points_receiver, games_server, games_receiver, sets_server, sets_receiver)
									}
								}
							}
						}
					}
				}
			}
			//Adiciona uma vitória no vencedor o game e exibe quem ganhou (server ou receiver)
			if set_winner_game == 0 {
				println("The game winner was the server")
				games_server++
			} else {
				println("The game winner was the receiver")
				games_receiver++
				
			}

		}
		/*
			Verifica se o sacador tem o numero de games suficientes para vencer o set.
			Caso afirmativo, adiona um set vencido na pontuação do sacador e exibe ele
			como vitorioso do set. Após isso, verifica se com adiação ele possui pontos 
			suficientes para vencer o match, caso seja, define-o como vencedor.

			O else funciona do mesmo modo, substituindo o server pelo receiver
		*/
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

	//Verifica o vencedor do game e exibe na tela
	if(winner_match == 0){
		println("The match winner was the server. THE END...")
	} else {
		println("The match winner was the receiver. THE END...")
	}
	println("+---------------------  THE END  ------------------------+")
	println("|               Server ", sets_server, " X ", sets_receiver," Receiver                |")
	println("+--------------------------------------------------------+")

}
