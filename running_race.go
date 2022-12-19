package main
import (
    "fmt"
    "time"
)
func emitter(out chan string) {
	message := "Entregou"
	fmt.Println("Message sent:", message)
	out <- message
}

func receiver(in chan string) {
	message := <- in
	fmt.Println("Message received:", message)
}
func main() {
    for i := 0; i <= 4; i++ {
		done := make(chan string)
		time.Sleep(1000 * time.Millisecond)
		go emitter(done)
		go receiver(done)
	}
}