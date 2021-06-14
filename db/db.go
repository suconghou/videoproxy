package db

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/suconghou/videoproxy/util"
)

var (
	db      *sql.DB
	dsn     = fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8&timeout=2s", os.Getenv("DB_AUTH"), os.Getenv("DB_ADDR"), os.Getenv("DB_NAME"))
	baseURL = os.Getenv("BASE_URL") // http://domain/video
	client  = &http.Client{
		Timeout: time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

type tableName string

const (
	TABLE_WHITELIST tableName = "whitelist"
	TABLE_CAPTIONS  tableName = "captions"
	TABLE_CACHEJSON tableName = "cachejson"
	TABLE_CACHEMPD  tableName = "cachempd"
)

func init() {
	if !strings.HasPrefix(dsn, "@tcp") {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			util.Log.Fatal(err)
		}
	}
}

// FindId 查某表中是否存在此ID
func FindId(id string, t tableName) (string, bool, error) {
	if db == nil {
		if baseURL == "" { // 不设置数据库,也不设置上游白名单,则全部放行
			return id, true, nil
		}
		if allowVideo(id) { // 使用上游白名单
			return id, true, nil
		}
		return "", false, nil
	}
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

// 请求上游json信息接口,这个是上游数据库会缓存的
// http 204 必然是不在白名单的, 200 是在白名单且能正常解析的, 500是在白名单有可能解析出错
func allowVideo(vid string) bool {
	var (
		u      = strings.Split(baseURL, ";")
		status int
		err    error
	)
	for _, part := range u {
		if part == "" {
			continue
		}
		var target = fmt.Sprintf("%s/%s.json", part, vid)
		status, err = httpStatus(target)
		if err != nil {
			util.Log.Print(err)
			continue
		}
		if status == http.StatusNoContent {
			return false
		}
		if status == http.StatusOK {
			return true
		}
	}
	if status == http.StatusInternalServerError {
		return true
	}
	return false
}

func httpStatus(target string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	return res.StatusCode, nil
}

func FindCaption(id string, lang string) (string, bool, error) {
	if db == nil {
		return "", false, nil
	}
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
	if db == nil {
		return nil
	}
	stmt, err := db.Prepare(fmt.Sprintf("REPLACE INTO %s (`id`, `lang`, `data`, `time`) VALUES (?, ?, ?, %d)", TABLE_CAPTIONS, time.Now().Unix()))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, lang, data)
	return err
}

func GetCacheItem(id string, table tableName) (string, bool, error) {
	if db == nil {
		return "", false, nil
	}
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
	if db == nil {
		return nil
	}
	stmt, err := db.Prepare(fmt.Sprintf("REPLACE INTO %s (`id`, `data`, `time`) VALUES (?, ?, %d)", table, time.Now().Unix()))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, data)
	return err
}
