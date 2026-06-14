package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	questions := make(map[string]string)

	problemsFilePath := flag.String("p", "./problems.csv", "Path to the CSV file containing questions and answers")
	timeLimit := flag.Int("t", 30, "Time limit for the quiz in seconds")
	flag.Parse()

	content, err := os.ReadFile(*problemsFilePath)
	if err != nil {
		fmt.Printf("Failed to read the CSV file: %s", err)
		return
	}

	r := csv.NewReader(strings.NewReader(string(content)))
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Failed to read the CSV line: %s", err)
			return
		}

		questions[strings.TrimSpace(record[0])] = strings.TrimSpace(record[1])
	}

	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

loop:
	for question, answer := range questions {
		fmt.Printf("\nQuestion: %s\nAnswer: ", question)

		answerCh := make(chan string)
		go func() {
			var response string
			fmt.Scanf("%s\n", &response)
			answerCh <- response
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			break loop
		case response := <-answerCh:
			if response == answer {
				correct++
				fmt.Println("Correct!")
			} else {
				fmt.Printf("Wrong! The correct answer is: %s\n", answer)
			}
		}
	}

	fmt.Printf("\nYou got %d out of %d questions correct.", correct, len(questions))
}
