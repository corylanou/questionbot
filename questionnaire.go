package questionBot

import (
	"errors"
	"strconv"
)

type PrefixType int

const (
	Alpha PrefixType = iota
	Numeric
)

type Questionnaire struct {
	Questions  Questions
	PrefixType PrefixType
	index      int
}

func NewQuestionnaire(questions Questions) *Questionnaire {
	return &Questionnaire{
		index:      -1,
		Questions:  questions,
		PrefixType: Alpha,
	}
}

func (q *Questionnaire) Next() (*Question, int) {
	if q.index < len(q.Questions) {
		q.index++
	}
	if len(q.Questions) > q.index {
		return q.Questions[q.index], q.index
	}
	return nil, -1
}

func (q *Questionnaire) Back() (*Question, int) {
	if q.index != -1 {
		q.index--
	}
	if q.index > -1 {
		return q.Questions[q.index], q.index
	}
	return nil, -1
}

func (q *Questionnaire) Answer(choice string) error {
	if q.index < 0 && q.index < len(q.Questions) {
		return errors.New(ErrNoQuestionSelected)
	}

	i, err := q.PrefixToInt(choice)
	if err != nil {
		return err
	}
	if err := q.Questions[q.index].SelectAnswer(i); err != nil {
		return err
	}

	return nil
}

func (q *Questionnaire) IntToPrefix(i int) string {
	const (
		a = 'a'
		z = 'z'
	)

	switch q.PrefixType {
	case Alpha:
		if i > z-a {
			panic("too many choices")
		}
		return string(a + i)
	case Numeric:
		return strconv.Itoa(i)
	default:
		panic("unreachable code")
	}

}

func (q *Questionnaire) PrefixToInt(prefix string) (int, error) {
	const (
		a = 'a'
		z = 'z'
		A = 'A'
		Z = 'Z'
	)

	if q.PrefixType == Alpha {
		if len(prefix) > 1 || prefix == "" {
			return -1, errors.New(ErrInvalidChoice)
		}
		p := prefix[0]

		if (p < a || p > z) && (p < A || p > Z) {
			return -1, errors.New(ErrInvalidChoice)
		}
		if p >= a && p <= z {
			return int(p - a), nil
		}
		if p >= A && p <= Z {
			return int(p - A), nil
		}
		return -1, errors.New(ErrInvalidChoice)
	}

	i, e := strconv.Atoi(string(prefix))
	if i < 0 {
		return -1, errors.New(ErrInvalidChoice)
	}
	return i, e
}

func (q *Questionnaire) Completed() bool {
	for _, q := range q.Questions {
		if !q.Completed() {
			return false
		}
	}
	return true
}
