package leetcode

import (
	"os"
	"path/filepath"
)

type Submission struct {
	Language  string `json:"lang"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status_display"`
	Code      string `json:"code"`
}

func (s Submission) WriteTo(fileName string) error {
	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	if s.Language == "golang" {
		_, err = f.WriteString("package main\n")
		if err != nil {
			return err
		}
	}
	_, err = f.WriteString(s.Code)
	if err != nil {
		return err
	}
	return nil
}

func (s Submission) IsAccepted() bool {
	return s.Status == "Accepted"
}

type SubmissionsResponse struct {
	SubmissionsDump []Submission `json:"submissions_dump"`
}

func (r SubmissionsResponse) GetLatestAcceptedSubmissions(languages map[string]string) map[string]Submission {
	subs := map[string]Submission{}
	for _, sub := range r.SubmissionsDump {
		if !sub.IsAccepted() {
			continue
		}
		lang, ok := languages[sub.Language]
		if !ok {
			continue
		}
		s, ok := subs[lang]
		if !ok {
			subs[lang] = sub
			continue
		}
		if s.Timestamp < sub.Timestamp {
			subs[lang] = sub
		}
	}
	return subs
}
