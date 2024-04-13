package database

import (
	"database/sql"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/snippetaccumulator/configloader"
	"github.com/snippetaccumulator/snac/internal/backend/model"
)

func connect(t *testing.T) Database {
	mockLoader := configloader.NewMockLoader(map[string]interface{}{
		"Database.Name": "snac-test",
	})

	// Url and AuthToken are sensitive, so loaded via .env.test in this directory
	testEnvData, err := os.ReadFile("./.env.test")
	if err != nil {
		t.Errorf("Error reading .env.test: %v", err)
	}

	testEnv := make(map[string]string)

	for _, line := range strings.Split(string(testEnvData), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		testEnv[key] = value
	}

	for _, entry := range []string{"Database.Url", "Database.AuthToken"} {
		if _, ok := testEnv[entry]; !ok {
			t.Errorf("missing %s in .env.test", entry)
		} else {
			mockLoader.Override(entry, testEnv[entry])
		}
	}

	db, err := NewDB(mockLoader)
	if err != nil {
		t.Errorf("NewDB() error = %v", err)
	}

	if db == nil {
		t.Errorf("NewDB() db = nil, want not nil")
	}
	if db.DB == nil {
		t.Errorf("NewDB() db.DB = nil, want not nil")
	}
	if db.connector == nil {
		t.Errorf("NewDB() db.connector = nil, want not nil")
	}
	if db.dir == "" {
		t.Errorf("NewDB() db.dir = \"\", want not \"\"")
	}

	err = db.Ping()
	if err != nil {
		t.Errorf("db.Ping() error = %v", err)
	}

	return db
}

func teamCreateIfNotExist(connection Database, t *testing.T) string {
	teamName := "test-teamA"

	row := connection.(*DB).QueryRow("SELECT name FROM teams where name = ?", teamName)
	var gotName string
	err := row.Scan(&gotName)
	if err != nil && err != sql.ErrNoRows {
		t.Errorf("Error checking if team-testA exists: %+v", err)
	}
	if teamName == gotName {
		return teamName
	}

	err = connection.InsertTeam(teamName, teamName, "password", "password")
	if err != nil {
		t.Errorf("Error inserting team-testA: %+v", err)
	}

	return teamName
}

func insert(snippet model.Snippet, connection Database, t *testing.T) {
	_, err := connection.InsertSnippet(snippet)
	if err != nil {
		t.Errorf("Error inserting snippet: %v", err)
	}
}

func deleteAll(connection Database, t *testing.T) {
	_, err := connection.(*DB).Query("DELETE FROM snippets")
	if err != nil {
		t.Errorf("Error emptying snippet table: %+v", err)
	}
}

func TestInsertSnippet(t *testing.T) {
	connection := connect(t)
	defer connection.Close()
	teamName := teamCreateIfNotExist(connection, t)

	snippet := model.NewSnippetBuilder("test1", teamName).Build()
	insert(snippet, connection, t)
}

func TestGetSnippet(t *testing.T) {
	connection := connect(t)
	defer connection.Close()
	teamName := teamCreateIfNotExist(connection, t)

	putSnippet := model.NewSnippetBuilder("test1", teamName).
		WithContent("content").
		WithDescription("desc").
		WithLanguage("lang").
		WithTags([]string{"tag1", "tag2"}).
		Build()
	insert(putSnippet, connection, t)

	snippet, err := connection.GetByID(putSnippet.ID)
	if err != nil {
		t.Errorf("Error getting snippet: %v", err)
	}

	if snippet.ID != putSnippet.ID {
		t.Errorf("Got unexpected ID exp: %v, act: %v", putSnippet.ID, snippet.ID)
	}
	if snippet.TeamID != putSnippet.TeamID {
		t.Errorf("Got unexpected TeamID exp: %v, act: %v", putSnippet.TeamID, snippet.TeamID)
	}
	if snippet.Title != putSnippet.Title {
		t.Errorf("Got unexpected Title exp: %v, act: %v", putSnippet.Title, snippet.Title)
	}
	if snippet.Description != putSnippet.Description {
		t.Errorf("Got unexpected Description exp: %v, act: %v", putSnippet.Description, snippet.Description)
	}
	if snippet.Language != putSnippet.Language {
		t.Errorf("Got unexpected Language exp: %v, act: %v", putSnippet.Language, snippet.Language)
	}
	if reflect.DeepEqual(snippet.Tags, putSnippet.Tags) {
		t.Errorf("Got unexpected Tags exp: %v, act: %v", putSnippet.Tags, snippet.Tags)
	}
}

func TestGetByTeam(t *testing.T) {
	connection := connect(t)
	defer connection.Close()
	deleteAll(connection, t)
	teamName := teamCreateIfNotExist(connection, t)

	snippet := model.NewSnippetBuilder("test1", teamName).Build()
	insert(snippet, connection, t)

	partials, err := connection.GetByTeamID(teamName)
	if err != nil {
		t.Errorf("error getting partials: %+v", err)
	}

	if len(partials) != 1 {
		t.Errorf("Got wrong number of partials: exp %d, act: %d", 1, len(partials))
	}

	if partials[0].ID != snippet.ID {
		t.Errorf("Got unexpected ID exp: %v, act: %v", snippet.ID, partials[0].ID)
	}
	if partials[0].TeamID != snippet.TeamID {
		t.Errorf("Got unexpected TeamID exp: %v, act: %v", snippet.TeamID, partials[0].TeamID)
	}
	if partials[0].Title != snippet.Title {
		t.Errorf("Got unexpected Title exp: %v, act: %v", snippet.Title, partials[0].Title)
	}
	if reflect.DeepEqual(partials[0].Tags, snippet.Tags) {
		t.Errorf("Got unexpected Tags exp: %v, act: %v", snippet.Tags, partials[0].Tags)
	}
}

func TestCheckTeamPassword(t *testing.T) {
	connection := connect(t)
	defer connection.Close()
	deleteAll(connection, t)
	teamName := teamCreateIfNotExist(connection, t)

	// test correct password regular
	correct, err := connection.CheckTeamPassword(teamName, "password", false)
	if !correct || err != nil {
		t.Errorf("Got wrong password: exp: true, act: %v, err: %v", correct, err)
	}

	// test correct password admin
	correct, err = connection.CheckTeamPassword(teamName, "password", true)
	if !correct || err != nil {
		t.Errorf("Got wrong password: exp: true, act: %v, err: %v", correct, err)
	}

	// test incorrect password regular
	correct, err = connection.CheckTeamPassword(teamName, "wrong", false)
	if correct || err != nil {
		t.Errorf("Got wrong password: exp: false, act: %v, err: %v", correct, err)
	}

	// test incorrect password admin
	correct, err = connection.CheckTeamPassword(teamName, "wrong", true)
	if correct || err != nil {
		t.Errorf("Got wrong password: exp: false, act: %v, err: %v", correct, err)
	}
}
