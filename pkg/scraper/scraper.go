package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jiachengxu/leetdoad/pkg/config"
	"github.com/jiachengxu/leetdoad/pkg/leetcode"

	"github.com/rs/zerolog/log"
)

const (
	leetcodeURL               = "https://leetcode.com"
	leetcodeAlgorithmProblems = "/api/problems/algorithms/"
	leetcodeSubmissions       = "/api/submissions/"
)

type filePattern struct {
	Language       string
	QuestionName   string
	QuestionNumber int
	Difficulty     string
}

type scraper struct {
	client        http.Client
	config        *config.InConfig
	includeHeader bool
}

func NewScraper(client http.Client, config *config.InConfig, includeHeader bool) *scraper {
	return &scraper{
		client:        client,
		config:        config,
		includeHeader: includeHeader,
	}
}

func (s *scraper) Scrape() error {
	questions, err := s.scrapeSolvedQuestions()
	if err != nil {
		return err
	}
	qNum := len(questions)
	log.Info().Msgf("Started scraping submissions for %d solved questions", qNum)
	tmpl, err := template.New("pattern").Parse(s.config.Pattern)
	if err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	for len(questions) > 0 {
		sleepTime := time.Duration((rand.Float64() + 1) * float64(time.Second))
		time.Sleep(sleepTime)
		lastIdx := len(questions) - 1
		q := questions[lastIdx]
		submissions, err := s.scrapeLatestAcceptedSubmissionsForQuestion(q.QuestionTitleSlug)
		if err != nil {
			log.Debug().
				Str("progress", fmt.Sprintf("[%d/%d]", qNum-lastIdx, qNum)).
				Int("question id", q.FrontendQuestionID).
				Str("reason", err.Error()).
				Msg("failed to scraper submissions, retrying")
			continue
		}
		for _, sub := range submissions {
			filePattern := filePattern{
				Language:       sub.Language,
				QuestionName:   q.QuestionTitleSlug,
				QuestionNumber: q.FrontendQuestionID,
				Difficulty:     q.Difficulty,
			}
			if err := tmpl.Execute(buf, filePattern); err != nil {
				return err
			}
			fileName := fmt.Sprintf("%s/%s.%s", pwd, buf.String(), s.config.LanguageMap[sub.Language])
			buf.Reset()
			if err := sub.WriteTo(fileName, q, s.includeHeader); err != nil {
				return err
			}
			log.Debug().
				Str("progress", fmt.Sprintf("[%d/%d]", qNum-lastIdx, qNum)).
				Int("question id", q.FrontendQuestionID).
				Str("file name", fileName).
				Msg("succeed")
		}
		log.Info().
			Str("progress", fmt.Sprintf("[%d/%d]", qNum-lastIdx, qNum)).
			Int("question id", q.FrontendQuestionID).
			Msg("succeed")
		questions = questions[:lastIdx]
	}
	return nil
}

func (s *scraper) scrapeSolvedQuestions() ([]leetcode.Question, error) {
	request, err := http.NewRequest(http.MethodGet, leetcodeURL+leetcodeAlgorithmProblems, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("cookie", s.config.Cookie)
	resp, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	stats := &leetcode.Statistics{}
	err = json.Unmarshal(responseBody, stats)
	if err != nil {
		return nil, err
	}
	return stats.GetSolvedQuestions(), nil
}

func (s *scraper) scrapeLatestAcceptedSubmissionsForQuestion(titleSlug string) (map[string]leetcode.Submission, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s/", leetcodeURL, leetcodeSubmissions, titleSlug), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("cookie", s.config.Cookie)

	resp, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to scrape submission for question: %s, status: %s", titleSlug, resp.Status)
	}
	submissionsResponse := &leetcode.SubmissionsResponse{}
	if err := json.Unmarshal(b, submissionsResponse); err != nil {
		return nil, err
	}
	return submissionsResponse.GetLatestAcceptedSubmissions(s.config.LanguageMap), nil
}
