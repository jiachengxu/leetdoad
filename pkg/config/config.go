package config

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var languageMap = map[string]string{
	"bash":       "sh",
	"c":          "c",
	"cpp":        "cpp",
	"golang":     "go",
	"java":       "java",
	"javascript": "js",
	"typescript": "ts",
	"python":     "py",
	"python3":    "py",
	"rust":       "rs",
	"ruby":       "rb",
	"scala":      "scala",
	"swift":      "swift",
}

type InConfig struct {
	Cookie      string
	LanguageMap map[string]string
	Pattern     string
}

func GetConfig(configFile, cookie string) (*InConfig, error) {
	if cookie == "" {
		cookie = os.Getenv("LEETCODE_COOKIE")
		if cookie == "" {
			return nil, fmt.Errorf("leetcode cookie cannot be empty, you must either pass it from --cookie flag or set LEETCODE_COOKIE env")
		}
	}
	c, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}
	if len(c.Languages) == 0 {
		return nil, fmt.Errorf("languages can not be empty")
	}
	availableLanguages := make([]string, len(languageMap))
	i := 0
	for lang := range languageMap {
		availableLanguages[i] = lang
		i++
	}
	userLanguages := map[string]string{}
	for _, lang := range c.Languages {
		if lang == "*" {
			userLanguages = languageMap
			break
		}
		if _, ok := languageMap[lang]; !ok {
			return nil, fmt.Errorf("cannot parse language type: %s, available languages: %s", lang, strings.Join(availableLanguages, ","))
		}
		userLanguages[lang] = languageMap[lang]
	}
	if _, err = template.New("pattern").Parse(c.Pattern); err != nil {
		return nil, err
	}
	return &InConfig{
		Cookie:      cookie,
		LanguageMap: userLanguages,
		Pattern:     c.Pattern,
	}, nil
}

type config struct {
	Languages []string `yaml:"languages"`
	Pattern   string   `yaml:"pattern"`
}

func loadConfig(configFile string) (*config, error) {
	abs, err := filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}
	fileInfo, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("the config file that you specified is a directory, not a file")
	}
	f, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c := &config{}
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}
