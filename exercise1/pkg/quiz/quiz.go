package quiz

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	CorrectAnswerCount   int
	IncorrectAnswerCount int
	Err                  error
}

func (r Result) String() string {
	return fmt.Sprintf("Quiz status: %d correct, %d incorrect.\n",
		r.CorrectAnswerCount, r.IncorrectAnswerCount,
	)
}

func Exec(ctx context.Context, quizFilepath string) Result {
	fd, err := os.Open(quizFilepath)
	if err != nil {
		return Result{Err: fmt.Errorf("os.Open: %w", err)}
	}
	csvReader := csv.NewReader(fd)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return Result{Err: fmt.Errorf("csv ReadAll: %w", err)}
	}

	scanner := bufio.NewScanner(os.Stdin)
	var correctCount, incorrectCount int
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) != 2 {
			return Result{Err: fmt.Errorf("%w: want 2, got %d", csv.ErrFieldCount, len(row))}
		}
		question, answer := row[0], row[1]
		fmt.Println("Question: ", question)

		scanner.Scan()
		text := scanner.Text()
		if strings.TrimSpace(text) != answer {
			incorrectCount++
		} else {
			correctCount++
		}
	}

	return Result{
		CorrectAnswerCount:   correctCount,
		IncorrectAnswerCount: incorrectCount,
	}
}
