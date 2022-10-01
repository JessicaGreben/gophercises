package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessicagreben/gophercises/exercise1/pkg/quiz"
)

var (
	quizFileName = flag.String("file", "problems.csv", "path to the file containing a quiz")
	timerSeconds = flag.Int("timer", 30, "the time allowed to complete the quiz in seconds")
)

var usage = `Usage of %s:

	exercise1 [options...]

Example usage:
	$ exercise1
	$ exercise1 -help
	$ exercise1 -file quiz.csv
	$ exercise1 -timer 60

Flags:
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
	}
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("err os.Getwd %v\n", err)
		os.Exit(1)
	}
	quizFilepath := filepath.Join(dir, "quizzes", *quizFileName)

	q, err := quiz.NewQuiz(context.Background(), quizFilepath, *timerSeconds, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Printf("err NewQuiz %v\n", err)
		os.Exit(1)
	}
	if err := q.Exec(); err != nil {
		fmt.Printf("err quiz.Exec %v\n", err)
		os.Exit(1)
	}
	fmt.Println(q.Result())
}
