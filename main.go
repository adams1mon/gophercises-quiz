package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}


func main() {
	// path := "./problems.csv"

	path := flag.String("path", "./problems.csv", "CSV file for the quiz questions")
	allowedTime := flag.Duration("time", 30 * time.Second, "Number of seconds allowed to answer a question")
	flag.Parse()

	f, e := os.Open(*path)
	check(e)
	defer f.Close()

	r := csv.NewReader(f)

	questionCount := 0
	correctCount := 0

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	go func() {
		<- sigchan
		os.Exit(0)
	}()

	for {

		q, e := r.Read()

		if e == io.EOF {
			break
		}
		check(e)

		questionCount++
		correct, e := strconv.Atoi(q[1])
		check(e)

		fmt.Printf("question: %s = ?\n", q[0])

		// start timer
		closeGame := func() {
			fmt.Println("game over, you didn't answer in time")
			fmt.Printf("you got %d correct out of %d questions.", correctCount, questionCount)
			os.Exit(0)
		}
		timer := time.AfterFunc(*allowedTime, closeGame)

		var ans int
		fmt.Scanf("%d\n", &ans)

		if ans == correct {
			correctCount++
		}
		timer.Stop()
	}

	fmt.Println("game ends")
	fmt.Printf("you got %d correct out of %d questions.", correctCount, questionCount)
}

