package model

type SnippetBuilder struct {
	snippet Snippet
}

func NewSnippetBuilder(title, teamId string) *SnippetBuilder {
	return &SnippetBuilder{
		snippet: Snippet{
			Title:  title,
			TeamID: teamId,
			ID:     NewID(),
		},
	}
}

func (b *SnippetBuilder) WithDescription(description string) *SnippetBuilder {
	b.snippet.Description = description
	return b
}

func (b *SnippetBuilder) WithTags(tags []string) *SnippetBuilder {
	b.snippet.Tags = tags
	return b
}

func (b *SnippetBuilder) WithLanguage(language string) *SnippetBuilder {
	b.snippet.Language = language
	return b
}

func (b *SnippetBuilder) WithContent(content string) *SnippetBuilder {
	b.snippet.Content = content
	return b
}

func (b *SnippetBuilder) Build() Snippet {
	b.snippet.ID = NewID()
	snippet := b.snippet
	b.Reset()
	return snippet
}

func (b *SnippetBuilder) Reset() {
	b.snippet = Snippet{}
}

type TeamBuilder struct {
	team Team
}
