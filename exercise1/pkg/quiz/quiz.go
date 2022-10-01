package quiz

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	ErrTimeout        = errors.New("time limit exceeded")
	ErrQuizIncomplete = errors.New("not all questions are answered")
)

type Quiz struct {
	ctx          context.Context
	timer        time.Duration
	questions    []question
	currQuestion int
	result       *Result
	w            io.Writer
	scanner      *bufio.Scanner
}

func NewQuiz(ctx context.Context, quizFilepath string, timer time.Duration, r io.Reader, w io.Writer) (*Quiz, error) {
	q := &Quiz{
		ctx:     ctx,
		timer:   timer,
		result:  &Result{},
		w:       w,
		scanner: bufio.NewScanner(r),
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
	qa := []question{}
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		qa = append(qa, newQuestion(row[0], row[1]))
	}
	q.questions = qa
	return nil
}

func (q *Quiz) NextQuestion() question {
	nextQuestion := q.questions[q.currQuestion]
	q.currQuestion++
	return nextQuestion
}

func (q *Quiz) Completed() bool {
	return q.currQuestion == len(q.questions)
}

func (q *Quiz) Result() Result {
	return *q.result
}

func (q *Quiz) UserInputAnswer() (string, error) {
	if ok := q.scanner.Scan(); !ok {
		return "", q.scanner.Err()
	}
	return q.scanner.Text(), nil
}

func (q *Quiz) Exec() error {
	ctx, cancel := context.WithTimeout(q.ctx, q.timer)
	defer cancel()

	quizCompleted := make(chan error, 1)
	go func() {
		for !q.Completed() {
			qa := q.NextQuestion()
			fmt.Fprintln(q.w, "Question: ", qa.question)
			usersAnswer, err := q.UserInputAnswer()
			if err != nil {
				quizCompleted <- err
			}
			if usersAnswer != qa.answer {
				q.result.IncorrectAnswerCount++
			} else {
				q.result.CorrectAnswerCount++
			}
		}
		quizCompleted <- nil
	}()

	select {
	case completed := <-quizCompleted:
		return completed
	case <-ctx.Done():
		return fmt.Errorf("%w %v", ErrTimeout, q.timer)
	}
}

func (q *Quiz) CorrectAnswerCount() int {
	return q.result.CorrectAnswerCount
}

func (q *Quiz) IncorrectAnswerCount() int {
	return q.result.IncorrectAnswerCount
}

type question struct {
	question string
	answer   string
}

func newQuestion(q, a string) question {
	return question{
		question: q,
		answer:   a,
	}
}

type Result struct {
	CorrectAnswerCount   int
	IncorrectAnswerCount int
}

func (r Result) String() string {
	return fmt.Sprintf("Quiz status: %d correct, %d incorrect.\n",
		r.CorrectAnswerCount, r.IncorrectAnswerCount,
	)
}
