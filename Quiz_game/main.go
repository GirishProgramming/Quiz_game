package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problemes.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "teh time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file: %s\n", *csvFilename))
	}
	defer file.Close() // Important: Close the file when done

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

problenloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problenloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}

		fmt.Printf("You Scored %d out of %d.\n", correct, len(problems))
	}
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, lines := range lines {
		ret[i] = problem{
			q: lines[0],
			a: lines[1],
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
