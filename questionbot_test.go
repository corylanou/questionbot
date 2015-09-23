package questionBot_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/corylanou/questionbot"
)

func Test_NextBack(t *testing.T) {
	qs := questionBot.Questions{
		{Text: "Question 1"},
		{Text: "Question 2"},
		{Text: "Question 3"},
	}

	qa := questionBot.NewQuestionnaire(qs)

	for j := 1; j <= 3; j++ {
		q, i := qa.Next()
		if exp, got := fmt.Sprintf("Question %d", j), q.Text; exp != got {
			t.Errorf("exp: %s, got: %s", exp, got)
		}
		if exp, got := j-1, i; exp != got {
			t.Errorf("exp: %d, got %d", exp, got)
		}
	}

	q, i := qa.Next()
	if exp, got := -1, i; exp != got {
		t.Errorf("exp: %d, got %d", exp, got)
	}
	if got := q; got != nil {
		t.Errorf("exp: nil, got: %s", got)
	}

	for j := 3; j >= 1; j-- {
		q, i := qa.Back()
		if exp, got := fmt.Sprintf("Question %d", j), q.Text; exp != got {
			t.Errorf("exp: %s, got: %s", exp, got)
		}
		if exp, got := j-1, i; exp != got {
			t.Errorf("exp: %d, got %d", exp, got)
		}
	}

	q, i = qa.Back()
	if exp, got := -1, i; exp != got {
		t.Errorf("exp: %d, got %d", exp, got)
	}
	if got := q; got != nil {
		t.Errorf("exp: nil, got: %s", got)
	}
}

func Test_QuestionnaireAnswer(t *testing.T) {
	q := &questionBot.Question{
		Choices: []string{
			"Blue",
			"Green",
			"Red",
		},
	}
	qa := questionBot.NewQuestionnaire(questionBot.Questions{q})

	err := qa.Answer("")
	if err == nil {
		t.Errorf("expected err, got nil")
	}
	err = qa.Answer("12")
	if err == nil {
		t.Errorf("expected err, got nil")
	}

	// Try to answer when no question is selected
	err = qa.Answer("b")
	if err == nil {
		t.Fatalf("expected error %v, got nil ", questionBot.ErrNoQuestionSelected)
	}

	nq, i := qa.Next()
	if nq == nil {
		t.Fatalf("no question found")
	}
	if i != 0 {
		t.Fatalf("unexpected question index")
	}

	// Answer when an answer is selected
	err = qa.Answer("b")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if qa.Questions[0].Answered != 1 {
		t.Fatalf("expected answer 1, got %d", qa.Questions[0].Answered)
	}

}

func Test_QuestionAnswer(t *testing.T) {
	q := questionBot.Question{
		Choices: []string{
			"Blue",
			"Green",
			"Red",
		},
	}

	// Make sure we can accept multiple answers for multple choice questions
	if err := q.SelectAnswer([]int{1, 2}...); err == nil {
		t.Errorf("exp error %s, got %v", questionBot.ErrTooManyChoices, err)
	} else if err.Error() != questionBot.ErrTooManyChoices {
		t.Errorf("exp error %s, got %v", questionBot.ErrTooManyChoices, err)
	}

	// Check that we can't answer a multiple choice with more than one
	q.MultipleChoice = false
	if err := q.SelectAnswer([]int{1, 2}...); err == nil {
		t.Errorf("exp err %s, got %v", questionBot.ErrTooManyChoices, err)
	} else if err.Error() != questionBot.ErrTooManyChoices {
		t.Errorf("exp err %s, got %v", questionBot.ErrTooManyChoices, err)
	}

	for i := -1; i <= len(q.Choices); i++ {
		err := q.SelectAnswer([]int{i}...)
		if i == -1 || i >= len(q.Choices) {
			if err == nil {
				t.Errorf("expected error for selecting choice %d, got %v", i, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for selecting choice %d, got %v", i, err)
		}
		if q.Answered != i {
			t.Errorf("expected answer to be %d, got %d", i, q.Answered)
		}
	}
}

func Test_Questionnaire_PrefixToInt(t *testing.T) {
	qa := questionBot.NewQuestionnaire(questionBot.Questions{})

	tests := []struct {
		prefix     string
		err        error
		i          int
		prefixType questionBot.PrefixType
	}{
		{
			prefix:     "a",
			i:          0,
			prefixType: questionBot.Alpha,
		},
		{
			prefix:     "A",
			i:          0,
			prefixType: questionBot.Alpha,
		},
		{
			prefix:     "!",
			err:        errors.New(questionBot.ErrInvalidChoice),
			i:          -1,
			prefixType: questionBot.Alpha,
		},
		{
			prefix:     "1",
			err:        errors.New(questionBot.ErrInvalidChoice),
			i:          -1,
			prefixType: questionBot.Alpha,
		},
		{
			prefix:     "0",
			i:          0,
			prefixType: questionBot.Numeric,
		},
		{
			prefix:     "1",
			i:          1,
			prefixType: questionBot.Numeric,
		},
		{
			prefix:     "55",
			i:          55,
			prefixType: questionBot.Numeric,
		},
		{
			prefix:     "-10",
			err:        errors.New(questionBot.ErrInvalidChoice),
			i:          -1,
			prefixType: questionBot.Numeric,
		},
		{
			prefix:     "-1",
			err:        errors.New(questionBot.ErrInvalidChoice),
			i:          -1,
			prefixType: questionBot.Numeric,
		},
	}

	for _, test := range tests {
		t.Logf("testing %q", test.prefix)
		qa.PrefixType = test.prefixType
		i, err := qa.PrefixToInt(test.prefix)
		if test.err == nil && err != nil {
			t.Errorf("exp: %v, got: %v", test.err, err)
		} else if test.err != nil && err == nil {
			t.Errorf("exp: %v, got: %v", test.err, err)
		}
		if i != test.i {
			t.Errorf("exp: %d, got: %d", test.i, i)
		}
	}
}
