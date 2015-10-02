package questionBot

import "io/ioutil"

type Config struct {
	PrefixType PrefixType
	DataPath   string
}

type Service struct {
	PrefixType     PrefixType
	Questionnaires Questionnaires
}

func NewService(c Config) (*Service, error) {
	var s Service
	s.PrefixType = c.PrefixType
	if c.PrefixType == "" {
		s.PrefixType = Alpha
	}

	var data string
	if b, e := ioutil.ReadFile(c.DataPath); e != nil {
		return nil, e
	} else {
		data = string(b)
	}
	if q, e := LoadQuestionnaires(data); e != nil {
		return nil, e
	} else {
		s.Questionnaires = q
	}

	for _, qa := range s.Questionnaires {
		qa.PrefixType = s.PrefixType
	}

	return &s, nil
}
