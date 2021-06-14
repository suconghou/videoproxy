package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/suconghou/videoproxy/util"
)

var (
	db  *sql.DB
	dsn = fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8&timeout=2s", os.Getenv("DB_AUTH"), os.Getenv("DB_ADDR"), os.Getenv("DB_NAME"))
)

type tableName string

const (
	TABLE_WHITELIST tableName = "whitelist"
	TABLE_CAPTIONS  tableName = "captions"
	TABLE_CACHEJSON tableName = "cachejson"
	TABLE_CACHEMPD  tableName = "cachempd"
)

func init() {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		util.Log.Fatal(err)
	}
}

// FindId 查某表中是否存在此ID
func FindId(id string, t tableName) (string, bool, error) {
	var retId string
	err := db.QueryRow(fmt.Sprintf("SELECT id FROM %s WHERE `id` = ? ", t), id).Scan(&retId)
	switch {
	case err == sql.ErrNoRows:
		return retId, false, nil
	case err != nil:
		return retId, false, err
	default:
		return retId, true, nil
	}
}

func FindCaption(id string, lang string) (string, bool, error) {
	var data string
	err := db.QueryRow(fmt.Sprintf("SELECT data FROM %s WHERE `id` = ? AND `lang` = ? AND time > 0", TABLE_CAPTIONS), id, lang).Scan(&data)
	switch {
	case err == sql.ErrNoRows:
		return data, false, nil
	case err != nil:
		return data, false, err
	default:
		return data, true, nil
	}
}

func SaveCaption(id string, lang string, data string) error {
	stmt, err := db.Prepare(fmt.Sprintf("REPLACE INTO %s (`id`, `lang`, `data`, `time`) VALUES (?, ?, ?, %d)", TABLE_CAPTIONS, time.Now().Unix()))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, lang, data)
	return err
}

func GetCacheItem(id string, table tableName) (string, bool, error) {
	var data string
	err := db.QueryRow(fmt.Sprintf("SELECT data FROM %s WHERE `id` = ? AND time > 0", table), id).Scan(&data)
	switch {
	case err == sql.ErrNoRows:
		return data, false, nil
	case err != nil:
		return data, false, err
	default:
		return data, true, nil
	}
}

func SaveCacheItem(id string, data string, table tableName) error {
	stmt, err := db.Prepare(fmt.Sprintf("REPLACE INTO %s (`id`, `data`, `time`) VALUES (?, ?, %d)", table, time.Now().Unix()))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, data)
	return err
}
