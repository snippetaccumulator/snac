package database

import "github.com/snippetaccumulator/snac/internal/backend/model"

type Database interface {
	GetByID(id model.ID) (model.Snippet, error)
	GetByTeamID(teamID string) ([]model.PartialSnippet, error)
	InsertSnippet(snippet model.Snippet) (model.Snippet, error)
	UpdateSnippet(snippet model.Snippet) error
	DeleteSnippet(id model.ID) error
	GetTeamByID(teamID string) (model.Team, error)
	InsertTeam(teamID string, displayName string, password string, adminPassword string) error
	UpdateTeam(team model.Team) error
	DeleteTeam(teamID string) error
	CheckTeamPassword(teamID string, password string, admin bool) (bool, error)
	Close()
}
