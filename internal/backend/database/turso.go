package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/snippetaccumulator/configloader"
	"github.com/snippetaccumulator/snac/internal/common"
	"github.com/tursodatabase/go-libsql"
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
	fmt.Println(dbCfg)

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
