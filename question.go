package questionBot

import "errors"

type Questions []*Question
type Question struct {
	Text           string
	Choices        []string
	MultipleChoice bool
	Skippable      bool
	OpenAnswer     bool
	Answered       []int
}

func NewQuestion(text string, choices []string, multipleChoice, skippable bool) *Question {
	q := Question{
		Text:           text,
		Choices:        choices,
		MultipleChoice: multipleChoice,
		Skippable:      skippable,
	}
	return &q
}

func (q *Question) SelectAnswer(choices ...int) error {
	if !q.MultipleChoice && len(choices) > 1 {
		return errors.New(ErrTooManyChoices)
	}

	// TODO make this work
	if len(choices) != 1 {
		return errors.New("multiple choices are not currently supported")
	}

	choice := choices[0]
	if choice < 0 || choice >= len(q.Choices) {
		return errors.New(ErrInvalidChoice)
	}

	// Store the valid choice
	q.Answered = choices

	return nil
}

func (q *Question) Completed() bool {
	return len(q.Answered) > 0
}
