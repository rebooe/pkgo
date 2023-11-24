package cache

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"time"

	_ "modernc.org/sqlite"
)

type sqlCache struct {
	db *sql.DB

	getStmt        *sql.Stmt
	setStmt        *sql.Stmt
	deleteStmt     *sql.Stmt
	autoDeleteStmt *sql.Stmt
}

func NewSqliteCache(db *sql.DB, table string) (*sqlCache, error) {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			"key" TEXT NOT NULL,
			"value" blob NOT NULL,
			"expire" integer NOT NULL,
			PRIMARY KEY ("key")
	  	);`, table),
	)
	if err != nil {
		return nil, err
	}

	getStmt, err := db.PrepareContext(context.TODO(), fmt.Sprintf("SELECT value, expire FROM %s WHERE key = ? AND expire >= ?", table))
	if err != nil {
		return nil, err
	}
	setStmt, err := db.PrepareContext(context.TODO(), fmt.Sprintf("INSERT OR REPLACE INTO %s (key, value, expire) VALUES (?, ?, ?)", table))
	if err != nil {
		return nil, err
	}
	deleteStmt, err := db.PrepareContext(context.TODO(), fmt.Sprintf("DELETE FROM %s WHERE key = ?", table))
	if err != nil {
		return nil, err
	}
	autoDeleteStmt, err := db.PrepareContext(context.TODO(), fmt.Sprintf("DELETE FROM %s WHERE expire < ?", table))
	if err != nil {
		return nil, err
	}

	return &sqlCache{
		db:             db,
		getStmt:        getStmt,
		setStmt:        setStmt,
		deleteStmt:     deleteStmt,
		autoDeleteStmt: autoDeleteStmt,
	}, nil
}

func (c *sqlCache) Get(key string, v any) (bool, time.Time, error) {
	var (
		value  []byte
		expire int64
	)
	err := c.getStmt.QueryRow(key, time.Now().Unix()).Scan(&value, &expire)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, time.Time{}, nil
		}
		return false, time.Time{}, err
	}

	buf := bytes.NewBuffer(value)
	decode := gob.NewDecoder(buf)
	if err := decode.Decode(v); err != nil {
		return false, time.Time{}, err
	}
	return true, time.Unix(expire, 0), nil
}

func (c *sqlCache) Set(key string, value any, t time.Duration) error {
	// 触发自动清理
	if rand.Intn(100) <= 2 {
		c.autoDelete()
	}

	buf := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(value); err != nil {
		return err
	}

	_, err := c.setStmt.Exec(key, buf.Bytes(), time.Now().Add(t).Unix())
	return err
}

func (c *sqlCache) Delete(key string) error {
	_, err := c.deleteStmt.Exec(key)
	return err
}

func (c *sqlCache) autoDelete() {
	c.autoDeleteStmt.Exec(time.Now().Unix())
}
