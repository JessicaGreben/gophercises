package quiz_test

import (
	//"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jessicagreben/gophercises/exercise1/pkg/quiz"
)

func newTestQuiz(t *testing.T, timer time.Duration, r io.Reader, w io.Writer) *quiz.Quiz {
	t.Helper()
	testQuizFilepath := filepath.Join("..", "..", "quizzes", "testdata.csv")
	q, err := quiz.NewQuiz(context.Background(), testQuizFilepath, timer, r, w)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	return q
}

type SleepyReader struct {
}

func (r SleepyReader) Read(p []byte) (n int, err error) {
	time.Sleep(time.Duration(50) * time.Millisecond)
	return 1, nil
}

func TestQuizExecTimeout(t *testing.T) {
	r := SleepyReader{}
	q := newTestQuiz(t, time.Duration(5)*time.Millisecond, r, io.Discard)
	if want, got := quiz.ErrTimeout, q.Exec(); !errors.Is(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestQuizExecCorrect(t *testing.T) {
	testCases := []struct {
		name           string
		input          []byte
		correctCount   int
		incorrectCount int
	}{
		{
			name:           "all correct",
			input:          []byte("A\nB\nC\n"),
			correctCount:   3,
			incorrectCount: 0,
		},
		{
			name:           "partial correct",
			input:          []byte("A\n2\nD\n"),
			correctCount:   1,
			incorrectCount: 2,
		},
		{
			name:           "all incorrect",
			input:          []byte("F\n2\nD\n"),
			correctCount:   0,
			incorrectCount: 3,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("%v\n", err)
			}
			if _, err = w.Write(tt.input); err != nil {
				t.Fatalf("%v\n", err)
			}
			w.Close()
			q := newTestQuiz(t, time.Duration(5)*time.Second, r, w)
			if got := q.Exec(); got != nil {
				t.Errorf("want: %v, got: %v", nil, got)
			}
			if want, got := true, q.Completed(); want != got {
				t.Errorf("want: %v got: %v", want, got)
			}
			if got, want := tt.correctCount, q.CorrectAnswerCount(); got != want {
				t.Errorf("want: %v, got: %v", want, got)
			}
			if got, want := tt.incorrectCount, q.IncorrectAnswerCount(); got != want {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}
