package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {

	csvFile := flag.String("file", "problems.csv", "Specify the CSV file containing problems")
	flag.Parse()
	duration := flag.Int("timer", 3, "Set the timer in secs")
	flag.Parse()

	read := bufio.NewReader(os.Stdin)
	fmt.Printf("Type in any key to start timer")
	fmt.Println()
	read.ReadString('\n')

	file, err := os.Open(*csvFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("failed to parse the file: %v\n", err)
	}

	var questions []string
	var answers []string

	for _, record := range data {
		if len(record) > 0 {
			question := record[0]
			questions = append(questions, question)

			length := len(record)
			answer := record[length-1]
			answers = append(answers, answer)
		}
	}

	// shuffle
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
		answers[i], answers[j] = answers[j], answers[i]
	})

	timer := time.NewTimer(time.Duration(*duration) * time.Second)
	var score int

	for i := 0; i < len(questions); i++ {

		fmt.Printf("Question %d: %s\n", i+1, questions[i])
		answerCh := make(chan string)
		go func() {
			var userAnswer string
			fmt.Scanf("%s\n", &userAnswer)
			answerCh <- userAnswer
		}()
		select {
		case <-timer.C:
			fmt.Printf("Your score %d, maximum possible score %d ", score, len(answers))
			return
		case userAnswer := <-answerCh:
			if userAnswer == answers[i] {
				score++
			}
		}
	}

	fmt.Printf("Your score %d, maximum possible score %d ", score, len(answers))
	fmt.Println()

}
