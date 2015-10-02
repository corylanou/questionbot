package questionBot_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/corylanou/questionbot"
	"github.com/davecgh/go-spew/spew"
)

func Test_Questionnaire_NextBack(t *testing.T) {
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

func Test_Questionnaire_Answer(t *testing.T) {
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

	if len(qa.Questions[0].Answered) != 1 && qa.Questions[0].Answered[0] != 1 {
		t.Fatalf("expected answer 1, got %d", qa.Questions[0].Answered)
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

func Test_Questionnaire_Load(t *testing.T) {
	data := `
[[questionnaire]]
	title="Ad"
	closing="Great! I'm so excited to build this with you. You should recieve an email bill and project timeline shortly. Let me know if you need anything else! I'm always here to help :)"
	[[questionnaire.question]]
		text="Neat! We're building an Ad. Could you tell me what the purpose of the ad is?"
		choices = [
			"Awareness. (I want as many new people to know about us as possible.)",
			"Direct Action. (I want to drive traffic to a specific product or to download my mobile app.)",
			"Lead Generation (I want to drive user or email signups)",
			"Customer Retention (I want to build loyalty among my current customers).",
		]
	[[questionnaire.question]]
		text="Is this an open question?"
`
	q, e := questionBot.LoadQuestionnaires(data)
	if e != nil {
		t.Fatal(e)
	}
	spew.Dump(q)
	if len(q) != 1 {
		t.Fatalf("wrong number of questionnaires, got %d, exp 1", len(q))
	}
	qa1 := q[0]
	if got, exp := len(qa1.Questions), 2; got != exp {
		t.Fatal("wrong number of questions, exp %d, got %d", exp, got)
	}
	q1, q2 := qa1.Questions[0], qa1.Questions[1]
	if got, exp := q1.Type, questionBot.SingleChoice; exp != got {
		t.Errorf("wrong question type.  exp %s, got %s", exp, got)
	}
	if exp, got := 4, len(q1.Choices); exp != got {
		t.Errorf("wrong number of choices. exp %d, got %d", exp, got)
	}

	if got, exp := q2.Type, questionBot.OpenQuestion; exp != got {
		t.Errorf("wrong question type.  exp %s, got %s", exp, got)
	}
	if exp, got := 0, len(q2.Choices); exp != got {
		t.Errorf("wrong number of choices. exp %d, got %d", exp, got)
	}

}
