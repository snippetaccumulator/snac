package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/snippetaccumulator/configloader"
	"github.com/snippetaccumulator/snac/internal/backend/model"
	"github.com/snippetaccumulator/snac/internal/common"
	"github.com/tursodatabase/go-libsql"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	dir       string
	connector *libsql.Connector
	*sql.DB
}

func (db *DB) Close() {
	os.RemoveAll(db.dir)
	db.connector.Close()
	db.DB.Close()
}

func NewDB(loader configloader.Loader) (*DB, error) {
	var commonConfig common.CommonConfig
	err := loader.Load(&commonConfig)
	if err != nil {
		return nil, err
	}

	dbCfg := commonConfig.Database

	dir, err := os.MkdirTemp("", "snac-*")
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dir, fmt.Sprintf("%s.db", dbCfg.Name))

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, dbCfg.Url, libsql.WithAuthToken(dbCfg.AuthToken))
	if err != nil {
		return nil, err
	}

	db := &DB{
		dir:       dir,
		connector: connector,
		DB:        sql.OpenDB(connector),
	}
	return db, nil
}

var fullSnippetSqlFields = "id, team_id, title, description, tags, language, content, last_modified"

func scanRowToDBSnippet(scanner interface {
	Scan(dest ...interface{}) error
}) (model.DBSnippet, error) {
	var dBSnippet model.DBSnippet
	err := scanner.Scan(&dBSnippet.ID, &dBSnippet.TeamID, &dBSnippet.Title, &dBSnippet.Description, &dBSnippet.Tags, &dBSnippet.Language, &dBSnippet.Content, &dBSnippet.LastModified)
	if err != nil {
		return model.DBSnippet{}, err
	}
	return dBSnippet, nil
}

func fullRowToSnippet(row *sql.Row) (model.Snippet, error) {
	dBSnippet, err := scanRowToDBSnippet(row)
	if err != nil {
		return model.Snippet{}, err
	}

	snippet, err := dBSnippet.ToSnippet()
	if err != nil {
		return model.Snippet{}, err
	}
	return snippet, nil
}

func fullRowsToSnippet(rows *sql.Rows) ([]model.Snippet, error) {
	var snippets []model.Snippet

	for rows.Next() {
		dBSnippet, err := scanRowToDBSnippet(rows)
		if err != nil {
			return nil, err
		}

		snippet, err := dBSnippet.ToSnippet()
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

var partialSnippetSqlFields = "id, team_id, title, tags"

func partialRowsToSnippets(rows *sql.Rows) ([]model.PartialSnippet, error) {
	var partialSnippets []model.PartialSnippet

	for rows.Next() {
		var dbPartialSnippet model.DBPartialSnippet
		err := rows.Scan(&dbPartialSnippet.ID, &dbPartialSnippet.TeamID, &dbPartialSnippet.Title, &dbPartialSnippet.Tags)
		if err != nil {
			return nil, err
		}

		partialSnippet := dbPartialSnippet.ToPartialSnippet()
		partialSnippets = append(partialSnippets, partialSnippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return partialSnippets, nil
}

func (db *DB) GetByID(id model.ID) (model.Snippet, error) {
	query := `SELECT ` + fullSnippetSqlFields + ` FROM snippets WHERE id = ?`
	row := db.QueryRow(query, id)

	return fullRowToSnippet(row)
}

func (db *DB) GetByTeamID(teamID string) ([]model.PartialSnippet, error) {
	query := `SELECT ` + partialSnippetSqlFields + ` FROM snippets WHERE team_id = ?`
	rows, err := db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return partialRowsToSnippets(rows)
}

func (db *DB) InsertSnippet(snippet model.Snippet) (model.Snippet, error) {
	snippet.LastModified = time.Now()
	dbSnippet := snippet.ToDBSnippet()
	query := `INSERT INTO snippets (id, team_id, title, description, tags, language, content, last_modified) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, dbSnippet.ID, dbSnippet.TeamID, dbSnippet.Title, dbSnippet.Description, dbSnippet.Tags, dbSnippet.Language, dbSnippet.Content, dbSnippet.LastModified)
	return snippet, err
}

func (db *DB) UpdateSnippet(snippet model.Snippet) error {
	query := `UPDATE snippets SET team_id = ?, title = ?, description = ?, tags = ?, language = ?, content = ?, last_modified = ? WHERE id = ?`
	_, err := db.Exec(query, snippet.TeamID, snippet.Title, snippet.Description, snippet.Tags, snippet.Language, snippet.Content, snippet.LastModified, snippet.ID)
	return err
}

func (db *DB) DeleteSnippet(id model.ID) error {
	query := `DELETE FROM snippets WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func (db *DB) GetTeamByID(teamID string) (model.Team, error) {
	query := `SELECT name, display_name, created, last_modified, password_hash, admin_hash FROM teams WHERE name = ?`
	row := db.QueryRow(query, teamID)
	var dbTeam model.DBTeam
	err := row.Scan(&dbTeam.Name, &dbTeam.DisplayName, &dbTeam.Created, &dbTeam.LastModified, &dbTeam.PasswordHash, &dbTeam.AdminHash)
	if err != nil {
		return model.Team{}, err
	}
	return dbTeam.ToTeam(), nil
}

func (db *DB) InsertTeam(teamId, displayName, password, adminPassword string) error {
	lastModified := time.Now().Format(time.RFC3339)
	created := lastModified
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedAdminPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO teams (name, display_name, created, last_modified, password_hash, admin_hash) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, teamId, displayName, created, lastModified, hashedPassword, hashedAdminPassword)
	return err
}

func (db *DB) UpdateTeam(team model.Team) error {
	team.LastModified = time.Now()
	dbTeam := team.ToDBTeam()
	query := `UPDATE teams SET display_name = ?, created = ?, last_modified = ?, password_hash = ?, admin_hash = ? WHERE name = ?`
	_, err := db.Exec(query, dbTeam.DisplayName, dbTeam.Created, dbTeam.LastModified, dbTeam.PasswordHash, dbTeam.AdminHash, dbTeam.Name)
	return err
}

func (db *DB) DeleteTeam(teamID string) error {
	query := `DELETE FROM teams WHERE name = ?`
	_, err := db.Exec(query, teamID)
	return err
}

func (db *DB) CheckTeamPassword(teamID string, password string, admin bool) (bool, error) {
	var hash_field string
	if admin {
		hash_field = "admin_hash"
	} else {
		hash_field = "password_hash"
	}

	var displayName string
	var hash string

	query := `SELECT display_name, ` + hash_field + ` FROM teams WHERE name = ?`
	row := db.QueryRow(query, teamID)
	err := row.Scan(&displayName, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("Team with name '%s' was not found", teamID)
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
