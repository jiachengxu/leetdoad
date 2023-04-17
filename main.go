package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/jiachengxu/leetdoad/pkg/config"
	"github.com/jiachengxu/leetdoad/pkg/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version   = "dev"
	commit    = ""
	date      = ""
	goVersion = ""
	website   = "https://github.com/jiachengxu/leetdoad"
)

func printVersion() {
	fmt.Printf(`leetdoad version: %s
commit: %s
built at: %s
go version: %s
%s
`, version, commit, date, goVersion, website)
}

type flags struct {
	configFilePath string
	cookie         string
	debug          bool
	version        bool
	header         bool
}

func main() {
	f := flags{}
	flag.StringVar(&f.configFilePath, "config-file", ".leetdoad.yaml", "Path of the leetdoad config file")
	flag.StringVar(&f.cookie, "cookie", "", "Cookie that used for scraping problems and solutions from Leetcode website, you can either pass it from here, or set LEETCODE_COOKIE env")
	flag.BoolVar(&f.debug, "debug", false, "Debug logs")
	flag.BoolVar(&f.version, "version", false, "Show the current leetdoad version")
	flag.BoolVar(&f.header, "header", false, "Add LeetCode VSCode extension header")
	flag.Parse()

	if f.version {
		printVersion()
		return
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if f.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	cfg, err := config.GetConfig(f.configFilePath, f.cookie)
	if err != nil {
		log.Fatal().Msgf("failed to get config: %s", err.Error())
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	if err := scraper.NewScraper(client, cfg, f.header).Scrape(); err != nil {
		log.Fatal().Msgf("failed to scrape solutions: %s", err.Error())
	}
}
