package quiz_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jessicagreben/gophercises/exercise1/pkg/quiz"
)

func newTestQuiz(t *testing.T, timer int, r io.Reader, w io.Writer) *quiz.Quiz {
	t.Helper()
	testQuizFilepath := filepath.Join("..", "..", "quizzes", "testdata.csv")
	q, err := quiz.NewQuiz(context.Background(), testQuizFilepath, timer, r, w)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	return q
}

func TestQuizInit(t *testing.T) {
	q := newTestQuiz(t, 5, nil, io.Discard)
	expectedQuestions := []string{"1", "2", "3"}
	expectedAnswers := []string{"A", "B", "C"}
	for i := 0; i < len(expectedQuestions); i++ {
		if want, got := false, q.Completed(); want != got {
			t.Errorf("want: %v got: %v", want, got)
		}
		if want, got := true, q.Next(); want != got {
			t.Errorf("want: %v, got: %v", want, got)
		}
		if want, got := expectedQuestions[i], q.Question(); want != got {
			t.Errorf("want: %s, got: %s", want, got)
		}
		if want, got := expectedAnswers[i], q.Answer(); want != got {
			t.Errorf("want: %s, got: %s", want, got)
		}
	}

	if want, got := false, q.Next(); want != got {
		t.Errorf("want: %v, got: %v", want, got)
	}
	if want, got := true, q.Completed(); want != got {
		t.Errorf("want: %v got: %v", want, got)
	}
}

func TestQuizExecTimeout(t *testing.T) {
	q := newTestQuiz(t, 0, ioutil.NopCloser(bytes.NewReader(nil)), io.Discard)
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
			q := newTestQuiz(t, 5, r, w)
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
