package database

import (
	"os"
	"strings"
	"testing"

	"github.com/snippetaccumulator/configloader"
)

func Test_NewDB(t *testing.T) {
	mockLoader := configloader.NewMockLoader(map[string]interface{}{
		"Database.Name": "snac-test",
	})

	// Url and AuthToken are sensitive, so loaded via .env.test in this directory
	testEnvData, err := os.ReadFile("./.env.test")
	if err != nil {
		t.Errorf("os.ReadFile() error = %v", err)
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
}
