package utils

import (
	"bytes"
	"time"

	"github.com/adrg/frontmatter"
)

const dateLayout = "2006-01-02T15:04:05Z"

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

// StripFrontMatter parses the YAML/TOML/JSON frontmatter from data into a typed
// FrontMatter and returns the remaining document body. When a date is present it
// populates both OrigDate (parsed time, for sorting) and Date (formatted for display).
// A parse error (including an invalid date) is returned so the caller can skip the file.
func StripFrontMatter(data []byte) ([]byte, FrontMatter, error) {
	var matter FrontMatter
	rest, err := frontmatter.Parse(bytes.NewReader(data), &matter)
	if err != nil {
		return nil, matter, err
	}
	if matter.Date != "" {
		matter.OrigDate, err = ParseDate(matter.Date)
		if err != nil {
			return nil, matter, err
		}
		matter.Date, err = ModifyDate(matter.Date)
		if err != nil {
			return nil, matter, err
		}
	}
	return rest, matter, nil
}

// ParseDate parses a frontmatter date in the canonical layout.
func ParseDate(date string) (time.Time, error) {
	return time.Parse(dateLayout, date)
}

// ModifyDate parses a frontmatter date and reformats it for display.
func ModifyDate(date string) (string, error) {
	t, err := ParseDate(date)
	if err != nil {
		return "", err
	}
	return t.Format("January 2, 2006 | 15:04"), nil
}
