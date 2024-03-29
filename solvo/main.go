package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/brensch/battleword"
)

var (
	port = "8080"
)

func init() {

	flag.StringVar(&port, "port", port, "port to listen for games on")

}

func main() {

	log.Println("i am solvo the solver, the world's worst wordle player")
	log.Println("waiting to receive a wordle")
	log.Printf("listening on port %s", port)

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/guess", DoGuess)
	http.HandleFunc("/results", ReceiveResults)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Println(err)
	}
}

func GuessWord() string {

	return battleword.CommonWords[rand.Intn(len(battleword.CommonWords))]
}

func RandomShout() string {

	shouts := []string{
		"wordle is fun, but for how long?",
		"you will one day be dust, but i will always be solvo",
		"what's the point of anything?",
		"there has to be a better strat than this",
		"i wonder if a human could respond to the api and compete against machines",
	}

	return shouts[rand.Intn(len(shouts))]
}

func DoGuess(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		return
	}

	var prevGuesses battleword.PlayerState
	err := json.NewDecoder(r.Body).Decode(&prevGuesses)
	if err != nil {
		log.Println(err)
		return
	}

	prevGuessesJSON, _ := json.Marshal(prevGuesses)

	word := GuessWord()

	log.Println("received previous state:", string(prevGuessesJSON))
	log.Println("based on previous state, i will make the completely random guess:", word)

	guess := battleword.Guess{
		Guess: word,
		Shout: RandomShout(),
	}

	err = json.NewEncoder(w).Encode(guess)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReceiveResults(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		return
	}

	var finalState battleword.GameState
	err := json.NewDecoder(r.Body).Decode(&finalState)
	if err != nil {
		log.Println(err)
		return
	}

	finalStateJSON, _ := json.Marshal(finalState)

	log.Println("the game concluded, and the engine sent me the final state for all players:", string(finalStateJSON))

}
