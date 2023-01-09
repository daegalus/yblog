package utils

import (
	"bytes"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/caarlos0/log"
)

type FrontMatter struct {
	Author   string    `yaml:"author"`
	OrigDate time.Time `yaml:"-"`
	Date     string    `yaml:"date"`
	Draft    bool      `yaml:"draft"`
	Slug     string    `yaml:"slug"`
	Title    string    `yaml:"title"`
	Type     string    `yaml:"type"`
	Tags     []string  `yaml:"tags"`
}

func StripFrontMatter(data []byte) ([]byte, FrontMatter) {
	var matter FrontMatter
	rest, err := frontmatter.Parse(bytes.NewReader(data), &matter)
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing frontmatter")
		panic(err)
	}
	matter.Date = ModifyDate(matter.Date)
	return rest, matter
}

func ModifyDate(date string) string {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing date")
		panic(err)
	}

	return t.Format("January 2, 2006 | 15:04")
}

func ParseDate(date string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing date")
		panic(err)
	}

	return t
}
