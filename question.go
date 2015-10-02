package questionBot

import "errors"

type QuestionType string

const (
	SingleChoice           QuestionType = "single-choice"
	OpenQuestion           QuestionType = "open"
	MultipleChoiceQuestion QuestionType = "multiple-choice"
)

type Questions []*Question
type Question struct {
	Text      string
	Choices   []string
	Type      QuestionType
	Skippable bool
	Answered  []int
}

func NewQuestion(text string, choices []string, questionType QuestionType, skippable bool) *Question {
	q := Question{
		Text:      text,
		Choices:   choices,
		Type:      questionType,
		Skippable: skippable,
	}
	return &q
}

func (q *Question) Init() {
	if q.Type == "" && len(q.Choices) == 0 {
		q.Type = OpenQuestion
	} else if q.Type == "" {
		q.Type = SingleChoice
	}
}

func (q *Question) MultipleChoice() bool {
	return q.Type == MultipleChoiceQuestion
}

func (q *Question) SelectAnswer(choices ...int) error {
	if !q.MultipleChoice() && len(choices) > 1 {
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
