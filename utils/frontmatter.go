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

func StripFrontMatter(data []byte) ([]byte, FrontMatter, error) {
	var matter FrontMatter
	rest, err := frontmatter.Parse(bytes.NewReader(data), &matter)
	if err != nil {
		log.WithField("error", err).Error("Error parsing frontmatter")
		return nil, matter, err
	}
	if matter.Date != "" {
		matter.Date, err = ModifyDate(matter.Date)
		if err != nil {
			return nil, matter, err
		}
	}
	return rest, matter, nil
}

func ModifyDate(date string) (string, error) {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.WithField("error", err).WithField("date", date).Error("Error modifying date")
		return "", err
	}

	return t.Format("January 2, 2006 | 15:04"), nil
}

func ParseDate(date string) (time.Time, error) {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.WithField("error", err).WithField("date", date).Error("Error parsing date")
		return time.Time{}, err
	}

	return t, nil
}
