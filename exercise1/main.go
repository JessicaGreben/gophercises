package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jessicagreben/gophercises/exercise1/pkg/quiz"
)

var quizFileName = flag.String("file", "problems.csv", "path to the file containing a quiz")
var timerSeconds = flag.Int("timer", 30, "the time allowed to complete quiz in seconds, default is 30s")

func main() {
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("err os.Getwd %v\n", err)
	}
	quizFilepath := filepath.Join(dir, "quizzes", *quizFileName)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timerSeconds)*time.Second)
	defer cancel()

	quizResult := make(chan quiz.Result, 1)
	go func() {
		quizResult <- quiz.Exec(ctx, quizFilepath)
	}()

	select {
	case result := <-quizResult:
		if result.Err != nil {
			fmt.Println(result.Err)
			return
		}
		fmt.Println(result)
		return
	case <-ctx.Done():
		fmt.Println("Exceeded quiz timelimit.")
		return
	}
}
