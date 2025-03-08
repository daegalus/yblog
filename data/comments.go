package data

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type DisqusPost struct {
	Feed             string          `json:"feed"`
	CleanTitle       string          `json:"clean_title"`
	Dislikes         int             `json:"dislikes"`
	Likes            int             `json:"likes"`
	Message          string          `json:"message"`
	RatingsEnabled   bool            `json:"ratingsEnabled"`
	IsSpam           bool            `json:"isSpam"`
	IsDeleted        bool            `json:"isDeleted"`
	Category         string          `json:"category"`
	AdsDisabled      bool            `json:"adsDisabled"`
	Author           string          `json:"author"`
	UserScore        int             `json:"userScore"`
	ID               string          `json:"id"`
	SignedLink       string          `json:"signedLink"`
	CreatedAt        string          `json:"createdAt"`
	HasStreaming     bool            `json:"hasStreaming"`
	RawMessage       string          `json:"raw_message"`
	IsClosed         bool            `json:"isClosed"`
	Link             string          `json:"link"`
	Slug             string          `json:"slug"`
	Forum            string          `json:"forum"`
	Identifiers      []string        `json:"identifiers"`
	Posts            int             `json:"posts"`
	UserSubscription bool            `json:"userSubscription"`
	ValidateAllPosts bool            `json:"validateAllPosts"`
	Title            string          `json:"title"`
	HighlightedPost  *string         `json:"highlightedPost"`
	Comments         []DisqusComment `json:"comments"`
}

type DisqusComment struct {
	EditableUntil string `json:"editableUntil"`
	Dislikes      int    `json:"dislikes"`
	NumReports    int    `json:"numReports"`
	Likes         int    `json:"likes"`
	Message       string `json:"message"`
	ID            string `json:"id"`
	CreatedAt     string `json:"createdAt"`
	Author        struct {
		Username                string `json:"username"`
		About                   string `json:"about"`
		Name                    string `json:"name"`
		Disable3rdPartyTrackers bool   `json:"disable3rdPartyTrackers"`
		IsPowerContributor      bool   `json:"isPowerContributor"`
		JoinedAt                string `json:"joinedAt"`
		ProfileUrl              string `json:"profileUrl"`
		Url                     string `json:"url"`
		Location                string `json:"location"`
		IsPrivate               bool   `json:"isPrivate"`
		SignedUrl               string `json:"signedUrl"`
		IsPrimary               bool   `json:"isPrimary"`
		IsAnonymous             bool   `json:"isAnonymous"`
		ID                      string `json:"id"`
		Avatar                  struct {
			Permalink string `json:"permalink"`
			Xlarge    struct {
				Permalink string `json:"permalink"`
				Cache     string `json:"cache"`
			} `json:"xlarge"`
			Cache string `json:"cache"`
			Large struct {
				Permalink string `json:"permalink"`
				Cache     string `json:"cache"`
			} `json:"large"`
			Small struct {
				Permalink string `json:"permalink"`
				Cache     string `json:"cache"`
			} `json:"small"`
			IsCustom bool `json:"isCustom"`
		} `json:"avatar"`
	} `json:"author"`
	Media                  []interface{} `json:"media"`
	IsSpam                 bool          `json:"isSpam"`
	IsDeletedByAuthor      bool          `json:"isDeletedByAuthor"`
	IsHighlighted          bool          `json:"isHighlighted"`
	Parent                 *int          `json:"parent"`
	IsApproved             bool          `json:"isApproved"`
	IsNewUserNeedsApproval bool          `json:"isNewUserNeedsApproval"`
	IsDeleted              bool          `json:"isDeleted"`
	IsFlagged              bool          `json:"isFlagged"`
	RawMessage             string        `json:"raw_message"`
	IsAtFlagLimit          bool          `json:"isAtFlagLimit"`
	CanVote                bool          `json:"canVote"`
	Forum                  string        `json:"forum"`
	Url                    string        `json:"url"`
	Points                 int           `json:"points"`
	ModerationLabels       []interface{} `json:"moderationLabels"`
	IsEdited               bool          `json:"isEdited"`
	Sb                     bool          `json:"sb"`
}

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

func GetLegacyComments() *map[string]*LegacyPost {
	filePath := "data/content/blog/legacy-comments.json"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	var legacyPosts map[string]*LegacyPost
	if err := json.Unmarshal(byteValue, &legacyPosts); err != nil {
		log.Fatalf("failed to unmarshal JSON: %s", err)
	}

	return &legacyPosts
}

func ProcessLegacyComments() *map[string]*LegacyPost {

	filePath := "data/disqus-comments-export.json"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	var site map[string]DisqusPost
	if err := json.Unmarshal(byteValue, &site); err != nil {
		log.Fatalf("failed to unmarshal JSON: %s", err)
	}

	legacyPosts := GenerateLegacyCommentsFromPosts(site)

	// write to file
	legacyFile, err := os.Create("data/content/blog/legacy-comments.json")
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	defer legacyFile.Close()

	legacyBytes, err := json.MarshalIndent(legacyPosts, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal JSON: %s", err)
	}

	_, err = legacyFile.Write(legacyBytes)
	if err != nil {
		log.Fatalf("failed to write file: %s", err)
	}

	return legacyPosts
}

func GenerateLegacyCommentsFromPosts(posts map[string]DisqusPost) *map[string]*LegacyPost {
	legacyPosts := &map[string]*LegacyPost{}
	allLegacyComments := &map[int]*LegacyComment{}

	for _, post := range posts {
		idInt, _ := strconv.Atoi(post.ID)
		legacyPost := LegacyPost{
			ID:       idInt,
			Title:    post.Title,
			Slug:     strings.ReplaceAll(strings.Replace(post.Slug, "_touch_devnull", "", 1), "_", "-"),
			Comments: &map[int]*LegacyComment{},
		}

		for _, comment := range post.Comments {
			legacyComment := CommentToLegacyComment(comment)
			(*allLegacyComments)[legacyComment.ID] = legacyComment
			if comment.Parent == nil {
				(*legacyPost.Comments)[legacyComment.ID] = legacyComment
			}
		}

		for _, comment := range *allLegacyComments {
			if comment.Parent > 0 {
				if slices.ContainsFunc((*allLegacyComments)[comment.Parent].Children, func(iComment *LegacyComment) bool {
					return comment.ID == iComment.ID
				}) {
					continue
				}
				(*allLegacyComments)[comment.Parent].Children = append((*allLegacyComments)[comment.Parent].Children, comment)
			}
		}
		(*legacyPosts)[legacyPost.Slug] = &legacyPost
	}
	return legacyPosts
}

func CommentToLegacyComment(comment DisqusComment) *LegacyComment {
	parent := -1
	if comment.Parent != nil {
		parent = *comment.Parent
	}
	idInt, _ := strconv.Atoi(comment.ID)
	legacyComment := LegacyComment{
		ID:             idInt,
		AuthorName:     comment.Author.Name,
		AuthorUsername: comment.Author.Username,
		AuthorLocation: comment.Author.Location,
		Message:        comment.Message,
		RawMessage:     comment.RawMessage,
		Likes:          comment.Likes,
		Dislikes:       comment.Dislikes,
		CreatedAt:      comment.CreatedAt,
		Parent:         parent,
		Children:       []*LegacyComment{},
	}
	return &legacyComment
}
