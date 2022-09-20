package quiz

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func Exec(quizFilepath string) error {
	fd, err := os.Open(quizFilepath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	csvReader := csv.NewReader(fd)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("csv ReadAll: %w", err)
	}
	var countCorrect, countIncorrect int
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) != 2 {
			return fmt.Errorf("%w: want 2, got %d", csv.ErrFieldCount, len(row))
		}
		question, answer := row[0], row[1]
		fmt.Println("Question: ", question)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		if strings.TrimSpace(text) != answer {
			countIncorrect++
		} else {
			countCorrect++
		}
	}
	fmt.Printf("Quiz completed. Total question count: %d. %d correct, %d incorrect.\n", len(rows)-1, countCorrect, countIncorrect)
	return nil
}
