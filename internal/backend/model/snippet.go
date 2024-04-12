package model

import (
	"math/rand"
	"strings"
	"time"
)

// ID is a unique identifier for a snippet. 5 Charactrers of uppercase letters and digits, except O and 0.
type ID string

func (id ID) String() string {
	return string(id)
}

// NewID generates a new ID for a snippet.
func NewID() ID {
	id := ""
	for i := 0; i < 5; i++ {
		id += string("ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"[rand.Intn(34)])
	}
	return ID(id)
}

// Snippet is a piece of code with a title, description, tags, language, and content.
type Snippet struct {
	ID           ID
	TeamID       string
	Title        string
	Description  string
	Tags         []string
	Language     string
	Content      string
	LastModified time.Time
}

type PartialSnippet struct {
	ID     ID
	TeamID string
	Title  string
	Tags   []string
}

// DBSnippet is a piece of code with a title, description, tags, language, and content, but with a string ID.
type DBSnippet struct {
	ID           string
	TeamID       string
	Title        string
	Description  string
	Tags         string
	Language     string
	Content      string
	LastModified string
}

const SnippetTableSql = `
CREATE TABLE IF NOT EXISTS snippets (
    id TEXT PRIMARY KEY,
    team_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    tags TEXT,
    language TEXT,
    content TEXT NOT NULL,
    last_modified TEXT NOT NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id)
);
`

type DBPartialSnippet struct {
	ID     string
	TeamID string
	Title  string
	Tags   string
}

func join(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return tags[0] + "," + join(tags[1:])
}

// ToDBSnippet converts a Snippet to a DBSnippet.
func (s Snippet) ToDBSnippet() DBSnippet {
	return DBSnippet{
		ID:           string(s.ID),
		TeamID:       s.TeamID,
		Title:        s.Title,
		Description:  s.Description,
		Tags:         join(s.Tags),
		Content:      s.Content,
		Language:     s.Language,
		LastModified: s.LastModified.Format(time.RFC3339),
	}
}

// ToSnippet converts a DBSnippet to a Snippet.
func (s DBSnippet) ToSnippet() (Snippet, error) {
	tags := []string{}
	if s.Tags != "" {
		tags = strings.Split(s.Tags, ",")
	} else {
		tags = []string{}
	}
	lastModified, err := time.Parse(time.RFC3339, s.LastModified)
	if err != nil {
		return Snippet{}, err
	}
	return Snippet{
		ID:           ID(s.ID),
		TeamID:       s.TeamID,
		Title:        s.Title,
		Description:  s.Description,
		Tags:         tags,
		Language:     s.Language,
		Content:      s.Content,
		LastModified: lastModified,
	}, nil
}

// ToDBPartialSnippet converts a PartialSnippet to a DBPartialSnippet.
func (s DBPartialSnippet) ToPartialSnippet() PartialSnippet {
	tags := []string{}
	if s.Tags != "" {
		tags = strings.Split(s.Tags, ",")
	} else {
		tags = []string{}
	}
	return PartialSnippet{
		ID:     ID(s.ID),
		TeamID: s.TeamID,
		Title:  s.Title,
		Tags:   tags,
	}
}
