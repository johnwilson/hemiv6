package main

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"strings"

	_ "github.com/johnwilson/bytengine/auth/mongo"
	_ "github.com/johnwilson/bytengine/bytestore/mongo"
	_ "github.com/johnwilson/bytengine/cmdhandler/base"
	_ "github.com/johnwilson/bytengine/datafilter/builtin"
	_ "github.com/johnwilson/bytengine/filesystem/mongo"
	_ "github.com/johnwilson/bytengine/parser/base"
	_ "github.com/johnwilson/bytengine/statestore/redis"
)

type MongoConfig struct {
	Addrs    []string `json:"addresses"`
	AuthDb   string   `json:"authdb"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Timeout  int64    `json:"timeout"`
}

type AuthConfig struct {
	*MongoConfig
	Plugin string `json:"plugin"`
	UserDb string `json:"userdb"`
}

type FileSystemConfig struct {
	*MongoConfig
	Plugin string `json:"plugin"`
}

type ByteStoreConfig struct {
	*MongoConfig
	Plugin  string `json:"plugin"`
	StoreDb string `json:"storedb"`
}

type StateConfig struct {
	Plugin   string `json:"plugin"`
	Address  string `json:"address"`
	Timeout  int64  `json:"timeout"`
	Password string `json:"password"`
	Database int64  `json:"database"`
}

type DataFilterConfig struct {
	Plugin string `json:"plugin"`
}

type ParserConfig struct {
	Plugin string `json:"plugin"`
}

type BytengineConfig struct {
	Authentication AuthConfig       `json:"authentication"`
	FileSystem     FileSystemConfig `json:"filesystem"`
	StateStore     StateConfig      `json:"statestore"`
	ByteStore      ByteStoreConfig  `json:"bytestore"`
	DataFilter     DataFilterConfig `json:"datafilter"`
	Parser         ParserConfig     `json:"parser"`
}

type TimeoutConfig struct {
	AuthToken    int64 `json:"authtoken"`
	UploadTicket int64 `json:"uploadticket"`
}

type ApplicationConfig struct {
	Workers int64         `json:"workers"`
	Port    int64         `json:"port"`
	Address string        `json:"address"`
	Timeout TimeoutConfig `json:"timeout"`
}

type Config struct {
	Bytengine json.RawMessage
	ApplicationConfig
}

func configToJSON() ([]byte, error) {
	var b []byte

	// client timeouts
	mt, err := strconv.ParseInt(os.Getenv("MONGO_TIMEOUT"), 10, 64)
	if err != nil {
		return b, err
	}
	rt, err := strconv.ParseInt(os.Getenv("REDIS_TIMEOUT"), 10, 64)
	if err != nil {
		return b, err
	}

	// parse mongodb connection url
	m_env := os.Getenv("MONGOLAB_URI")
	if len(m_env) == 0 {
		m_env = os.Getenv("MONGODB_URL")
	}
	println(m_env)
	m_url, err := url.Parse(m_env)
	if err != nil {
		return b, err
	}
	m_user := m_url.User.Username()
	m_pw, _ := m_url.User.Password()
	m_host := m_url.Host
	m_db := strings.TrimPrefix(m_url.Path, "/")

	mc := MongoConfig{
		[]string{m_host},
		m_db,
		m_user,
		m_pw,
		mt,
	}

	a_env := os.Getenv("AUTH_DB")
	if len(a_env) == 0 {
		a_env = m_db
	}
	ac := AuthConfig{
		&mc,
		"mongodb",
		a_env,
	}

	fc := FileSystemConfig{
		&mc,
		"mongodb",
	}

	b_env := os.Getenv("STORE_DB")
	if len(b_env) == 0 {
		b_env = m_db
	}
	bsc := ByteStoreConfig{
		&mc,
		"mongodb",
		b_env,
	}

	// parse redis connection url
	r_env := os.Getenv("REDISTOGO_URL")
	if len(r_env) == 0 {
		r_env = os.Getenv("REDIS_URL")
	}
	r_url, err := url.Parse(r_env)
	if err != nil {
		return b, err
	}
	r_pw, _ := r_url.User.Password()
	r_host := r_url.Host
	r_db, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
	if err != nil {
		return b, err
	}

	sc := StateConfig{
		"redis",
		r_host,
		rt,
		r_pw,
		r_db,
	}

	// get bytengine port
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		return b, err
	}

	// get bytengine port
	workers, err := strconv.ParseInt(os.Getenv("CMD_WORKERS"), 10, 64)
	if err != nil {
		return b, err
	}

	// build config struct
	dc := DataFilterConfig{"core"}
	pc := ParserConfig{"base"}
	bc := BytengineConfig{ac, fc, sc, bsc, dc, pc}
	tc := TimeoutConfig{60, 60}

	config := struct {
		Bytengine BytengineConfig `json:"bytengine"`
		Workers   int64           `json:"workers"`
		Port      int64           `json:"port"`
		Address   string          `json:"address"`
		Timeout   TimeoutConfig   `json:"timeout"`
	}{
		bc,
		workers,
		port,
		"",
		tc,
	}

	return json.Marshal(config)
}
