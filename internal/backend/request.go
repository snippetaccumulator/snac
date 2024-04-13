package backend

import (
	"fmt"

	"github.com/snippetaccumulator/snac/internal/backend/database"
	"github.com/snippetaccumulator/snac/internal/backend/model"
)

type Operation int

const (
	Get Operation = iota
	GetAllPartials
	Insert
	Update
	Delete
	InsertTeam
	UpdateTeam
	DeleteTeam
)

type Request struct {
	Operation Operation
	teamID    string
	password  string
	admin     bool

	Data any
}

type RequestBuilder struct {
	request Request
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		request: Request{},
	}
}

func (b *RequestBuilder) ForTeamByID(teamID string, password string, admin bool) *RequestBuilder {
	b.request.teamID = teamID
	b.request.password = password
	b.request.admin = admin
	return b
}

func (b *RequestBuilder) Get(snippetID model.ID) *RequestBuilder {
	b.request.Operation = Get
	b.request.Data = snippetID
	return b
}

func (b *RequestBuilder) GetAllPartials() *RequestBuilder {
	b.request.Operation = GetAllPartials
	return b
}

func (b *RequestBuilder) Insert(snippet model.Snippet) *RequestBuilder {
	b.request.Operation = Insert
	b.request.Data = snippet
	return b
}

func (b *RequestBuilder) Update(snippet model.Snippet) *RequestBuilder {
	b.request.Operation = Update
	b.request.Data = snippet
	return b
}

func (b *RequestBuilder) Delete(snippetID model.ID) *RequestBuilder {
	b.request.Operation = Delete
	b.request.Data = snippetID
	return b
}

func (b *RequestBuilder) NewTeam(team model.Team) *RequestBuilder {
	b.request.Operation = InsertTeam
	b.request.Data = team
	return b
}

func (b *RequestBuilder) UpdateTeam(team model.Team) *RequestBuilder {
	b.request.Operation = UpdateTeam
	b.request.Data = team
	return b
}

func (b *RequestBuilder) DeleteTeam(teamID string) *RequestBuilder {
	b.request.Operation = DeleteTeam
	b.request.Data = teamID
	return b
}

func (b *RequestBuilder) Build() Request {
	r := b.request
	b.Reset()
	return r
}

func (b *RequestBuilder) Reset() {
	b.request = Request{}
}

type RequestReturn int

const (
	ReturnSingleSnippet RequestReturn = iota
	ReturnSnippetList
	ReturnPartials
	ReturnTeam
	ReturnBoolean
	ReturnNone
)

func passwordCheckNeeded(op Operation) bool {
	switch op {
	case Get, GetAllPartials, Insert, Update, Delete, UpdateTeam, DeleteTeam:
		return true
	}
	return false
}

func (r Request) Execute(db database.Database) (any, RequestReturn, error) {
	if passwordCheckNeeded(r.Operation) {
		correctPassword, err := db.CheckTeamPassword(r.teamID, r.password, r.admin)
		if err != nil {
			return nil, ReturnNone, fmt.Errorf("Error while checking team password: %v", err)
		}
		if !correctPassword {
			return nil, ReturnNone, fmt.Errorf("Incorrect password for team '%s'", r.teamID)
		}
	}

	switch r.Operation {
	case GetAllPartials:
		partials, err := db.GetByTeamID(r.teamID)
		if err != nil {
			return nil, ReturnNone, fmt.Errorf("Error while executing GetAllPartials operation: %v", err)
		}
		return partials, ReturnPartials, nil
	case Get:
		id, ok := r.Data.(model.ID)
		if !ok {
			return nil, ReturnNone, fmt.Errorf("Request.Data for Get operation needs to be a string")
		}
		snippet, err := db.GetByID(id)
		if err != nil {
			return nil, ReturnNone, fmt.Errorf("Error while executing Get for '%s': %v", r.Data, err)
		}
		return snippet, ReturnSingleSnippet, nil
	case Insert:
		snippet, ok := r.Data.(model.Snippet)
		if !ok {
			return false, ReturnBoolean, fmt.Errorf("Request.Data for Insert operation needs to be a snippet")
		}
		snippet, err := db.InsertSnippet(snippet)
		if err != nil {
			return nil, ReturnNone, fmt.Errorf("Error while executing Insert for '%v': %v", r.Data, err)
		}
		return snippet, ReturnSingleSnippet, nil
	case Update:
		snippet, ok := r.Data.(model.Snippet)
		if !ok {
			return false, ReturnBoolean, fmt.Errorf("Request.Data for Update operation needs to be a snippet")
		}
		err := db.UpdateSnippet(snippet)
		if err != nil {
			return false, ReturnBoolean, fmt.Errorf("Error while executing Insert for '%v': %v", r.Data, err)
		}
		return true, ReturnBoolean, nil
	case Delete:
		snippetID, ok := r.Data.(model.ID)
		if !ok {
			return false, ReturnBoolean, fmt.Errorf("Request.Data for Delete operation needs to be a snippet")
		}
		err := db.DeleteSnippet(snippetID)
		if err != nil {
			return false, ReturnBoolean, fmt.Errorf("Error while executing Insert for '%v': %v", r.Data, err)
		}
		return true, ReturnBoolean, nil
	case InsertTeam:
		team, ok := r.Data.(model.Team)
		if !ok {
			return nil, ReturnNone, fmt.Errorf("Request.Data for InsertTeam operation needs to be a team")
		}
		err := db.InsertTeam(team.Name, team.DisplayName, team.PasswordHash, team.AdminHash)
		if err != nil {
			return nil, ReturnNone, fmt.Errorf("Error while executing InsertTeam for '%v': %v", r.Data, err)
		}
		return true, ReturnBoolean, nil
	case UpdateTeam:
		team, ok := r.Data.(model.Team)
		if !ok {
			return nil, ReturnNone, fmt.Errorf("Request.Data for UpdateTeam operation needs to be a team")
		}
		err := db.UpdateTeam(team)
		if err != nil {
			return nil, ReturnNone, fmt.Errorf("Error while executing UpdateTeam for '%v': %v", r.Data, err)
		}
		return true, ReturnBoolean, nil
	case DeleteTeam:
		teamId, ok := r.Data.(string)
		if !ok {
			return nil, ReturnNone, fmt.Errorf("Request.Data for DeleteTeam operation needs to be a string")
		}
		err := db.DeleteTeam(teamId)
		if err != nil {
			return nil, ReturnNone, fmt.Errorf("Error while executing DeleteTeam for '%v': %v", r.Data, err)
		}
		return nil, ReturnNone, nil
	}

	return nil, ReturnNone, nil
}
