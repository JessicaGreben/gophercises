package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessicagreben/gophercises/exercise1/pkg/quiz"
)

var customFileName = flag.String("file", "problems.csv", "path of a file containing a quiz")

func main() {
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("err os.Getwd %v\n", err)
	}
	quizFilepath := filepath.Join(dir, "quizzes", *customFileName)
	if err := quiz.Exec(quizFilepath); err != nil {
		fmt.Printf("err quiz.Exec %v\n", err)
	}
}
