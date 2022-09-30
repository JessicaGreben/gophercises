package quiz

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type Quiz struct {
	ctx             context.Context
	timer           time.Duration
	questionAnswers []questionAnswer
	currQuestion    int
	currAnswer      int
	result          *Result
}

func NewQuiz(ctx context.Context, quizFilepath string, timerSeconds int) (*Quiz, error) {
	q := &Quiz{
		ctx:    ctx,
		timer:  time.Duration(timerSeconds) * time.Second,
		result: &Result{},
	}
	err := q.parseQuestionAnswers(quizFilepath)
	return q, err
}

func (q *Quiz) parseQuestionAnswers(quizFilepath string) error {
	fd, err := os.Open(quizFilepath)
	if err != nil {
		return err
	}
	csvReader := csv.NewReader(fd)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	qa := []questionAnswer{}
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		qa = append(qa, newQuestionAnswer(row[0], row[1]))
	}
	q.questionAnswers = qa
	return nil
}

func (q *Quiz) Next() bool {
	return q.currQuestion < len(q.questionAnswers)
}

func (q *Quiz) Question() string {
	nextQuestion := q.questionAnswers[q.currQuestion]
	q.currQuestion++
	return nextQuestion.question
}

func (q *Quiz) Answer() string {
	nextAnswer := q.questionAnswers[q.currAnswer]
	q.currAnswer++
	return nextAnswer.answer
}

func (q *Quiz) Completed() bool {
	return q.currAnswer == len(q.questionAnswers)
}

func (q *Quiz) Result() Result {
	return *q.result
}

func (q *Quiz) Exec() error {
	ctx, cancel := context.WithTimeout(q.ctx, q.timer)
	defer cancel()

	quizCompleted := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for q.Next() {
			fmt.Println("Question: ", q.Question())
			if ok := scanner.Scan(); !ok {
				quizCompleted <- scanner.Err()
			}
			text := scanner.Text()
			if strings.TrimSpace(text) != q.Answer() {
				q.result.IncorrectAnswerCount++
			} else {
				q.result.CorrectAnswerCount++
			}
		}
		if q.Completed() {
			quizCompleted <- nil
		} else {
			quizCompleted <- fmt.Errorf("not all questions are answered")
		}
	}()

	select {
	case completed := <-quizCompleted:
		return completed
	case <-ctx.Done():
		return fmt.Errorf("exceeded quiz time limit %v", q.timer)
	}
}

type questionAnswer struct {
	question string
	answer   string
}

func newQuestionAnswer(q, a string) questionAnswer {
	return questionAnswer{
		question: q,
		answer:   a,
	}
}

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
