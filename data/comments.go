package data

import (
	"encoding/json"
)

type LegacyPost struct {
	ID       int                     `json:"id"`
	Title    string                  `json:"title"`
	Slug     string                  `json:"slug"`
	Comments *map[int]*LegacyComment `json:"comments"`
}
type LegacyComment struct {
	ID             int              `json:"id"`
	AuthorName     string           `json:"author_name"`
	AuthorUsername string           `json:"author_username"`
	AuthorLocation string           `json:"author_location"`
	Message        string           `json:"message"`
	RawMessage     string           `json:"raw_message"`
	Likes          int              `json:"likes"`
	Dislikes       int              `json:"dislikes"`
	CreatedAt      string           `json:"created_at"`
	Parent         int              `json:"parent"`
	Children       []*LegacyComment `json:"children"`
}

// GetLegacyComments loads the migrated Disqus comments from the embedded content.
// Comments are non-essential, so callers should treat an error as "no comments"
// rather than aborting the build.
func GetLegacyComments() (map[string]*LegacyPost, error) {
	raw, err := Content.ReadFile("content/blog/legacy-comments.json")
	if err != nil {
		return nil, err
	}
	var legacyPosts map[string]*LegacyPost
	if err := json.Unmarshal(raw, &legacyPosts); err != nil {
		return nil, err
	}
	return legacyPosts, nil
}
