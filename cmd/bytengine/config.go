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
	m_url, err := url.Parse(os.Getenv("MONGODB_URL"))
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

	ac := AuthConfig{
		&mc,
		"mongodb",
		os.Getenv("AUTH_DB"),
	}

	fc := FileSystemConfig{
		&mc,
		"mongodb",
	}

	bsc := ByteStoreConfig{
		&mc,
		"mongodb",
		os.Getenv("STORE_DB"),
	}

	// parse redis connection url
	r_url, err := url.Parse(os.Getenv("REDIS_URL"))
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
	workers, err := strconv.ParseInt(os.Getenv("WORKERS"), 10, 64)
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
