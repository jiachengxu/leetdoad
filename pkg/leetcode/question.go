package leetcode

var difficulties = map[int]string{
	1: "easy",
	2: "medium",
	3: "hard",
}

type Question struct {
	QuestionID         int    `json:"question_id"`
	FrontendQuestionID int    `json:"frontend_question_id"`
	QuestionTitleSlug  string `json:"question_title_slug"`
	Status             string `json:"status"`
	Difficulty         string `json:"difficulty"`
}

type Statistics struct {
	UserName        string `json:"user_name"`
	NumSolved       int    `json:"num_solved"`
	ACEasy          int    `json:"ac_easy"`
	ACMedium        int    `json:"ac_medium"`
	ACHard          int    `json:"ac_hard"`
	StatStatusPairs []struct {
		Stat struct {
			QuestionID         int    `json:"question_id"`
			FrontendQuestionID int    `json:"frontend_question_id"`
			QuestionTitleSlug  string `json:"question__title_slug"`
		} `json:"stat"`
		Status     string `json:"status"`
		Difficulty struct {
			Level int `json:"level"`
		} `json:"difficulty"`
	} `json:"stat_status_pairs"`
}

func (s Statistics) GetSolvedQuestions() []Question {
	var questions []Question
	for _, p := range s.StatStatusPairs {
		if p.Status != "ac" {
			continue
		}
		questions = append(questions, Question{
			QuestionID:         p.Stat.QuestionID,
			FrontendQuestionID: p.Stat.FrontendQuestionID,
			QuestionTitleSlug:  p.Stat.QuestionTitleSlug,
			Status:             p.Status,
			Difficulty:         difficulties[p.Difficulty.Level],
		})
	}
	return questions
}
