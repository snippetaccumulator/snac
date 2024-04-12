package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Team struct {
	Name         string
	DisplayName  string
	Created      time.Time
	LastModified time.Time
	PasswordHash string
	AdminHash    string
}

type DBTeam struct {
	Name         string
	DisplayName  string
	Created      string
	LastModified string
	PasswordHash string
	AdminHash    string
}

const TeamTableSql = `
CREATE TABLE IF NOT EXISTS teams (
	name TEXT PRIMARY KEY,
	display_name TEXT NOT NULL,
	created TEXT NOT NULL,
	last_modified TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	admin_hash TEXT NOT NULL
);
`

func (t Team) ToDBTeam() DBTeam {
	return DBTeam{
		Name:         t.Name,
		DisplayName:  t.DisplayName,
		Created:      t.Created.Format(time.RFC3339),
		LastModified: t.LastModified.Format(time.RFC3339),
		PasswordHash: t.PasswordHash,
		AdminHash:    t.AdminHash,
	}
}

func (t DBTeam) ToTeam() Team {
	created, _ := time.Parse(time.RFC3339, t.Created)
	lastModified, _ := time.Parse(time.RFC3339, t.LastModified)
	return Team{
		Name:         t.Name,
		DisplayName:  t.DisplayName,
		Created:      created,
		LastModified: lastModified,
		PasswordHash: t.PasswordHash,
		AdminHash:    t.AdminHash,
	}
}

func NewTeam(name, displayName, password, adminPassword string) (Team, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedAdminPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return Team{}, err
	}
	return Team{
		Name:         displayName,
		DisplayName:  displayName,
		Created:      time.Now(),
		LastModified: time.Now(),
		PasswordHash: string(hashedPassword),
		AdminHash:    string(hashedAdminPassword),
	}, nil
}
