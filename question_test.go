package questionBot_test

import (
	"testing"

	"github.com/corylanou/questionbot"
)

func Test_Question_Answer(t *testing.T) {
	q := questionBot.NewQuestion("", []string{"Blue", "Green", "Red"}, false, false)

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
		if len(q.Answered) != 1 && q.Answered[0] != i {
			t.Errorf("expected answer to be %d, got %d", i, q.Answered)
		}
	}
}
