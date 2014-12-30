package mongo

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/johnwilson/bytengine"
	_ "github.com/johnwilson/bytengine/bytestore/diskv"
	_ "github.com/johnwilson/bytengine/parser/base"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

const (
	BFS_CONFIG = `
    {
        "addresses":["localhost:27017"],
        "authdb":"",
        "username":"",
        "password":"",
        "timeout":60
    }`
	BSTORE_CONFIG = `
    {
        "rootdir":"/tmp/diskv_data",
        "cachesize": 1
    }`
)

func TestDatabaseManagement(t *testing.T) {
	// get bst plugin
	bstore, err := bytengine.NewByteStore("diskv", BSTORE_CONFIG)
	assert.Nil(t, err, "bst not created")
	// get bfs plugin
	mfs, err := bytengine.NewFileSystem("mongodb", BFS_CONFIG, &bstore)
	assert.Nil(t, err, "bfs not created")

	// Clear all
	_, err = mfs.ClearAll()
	assert.Nil(t, err, "clear all failed")

	// Create databases
	err = mfs.CreateDatabase("db1")
	assert.Nil(t, err, "db1 not created")
	err = mfs.CreateDatabase("db2")
	assert.Nil(t, err, "db2 not created")

	// List databases
	list, err := mfs.ListDatabase("")
	assert.Nil(t, err, "listing dbs failed")
	assert.Len(t, list, 2, "all databases not created")

	db1_found := false
	db2_found := false
	for _, db := range list {
		switch db {
		case "db1":
			db1_found = true
		case "db2":
			db2_found = true
		default:
			continue
		}
	}
	assert.True(t, db1_found, "db1 not created")
	assert.True(t, db2_found, "db2 not created")

	// Delete database
	err = mfs.DropDatabase("db2")
	assert.Nil(t, err, "db2 deletion failed")

	// check database list
	list, err = mfs.ListDatabase("")
	assert.Nil(t, err, "listing dbs failed")
	assert.Len(t, list, 1, "database listing wrong")

	db1_found = false
	db2_found = false
	for _, db := range list {
		switch db {
		case "db1":
			db1_found = true
		case "db2":
			db2_found = true
		default:
			continue
		}
	}
	assert.True(t, db1_found, "db1 not found after bytengine.dropdatabase")
	assert.False(t, db2_found, "db2 not deleted from bfs")
}

func TestContentManagement(t *testing.T) {
	// get bst plugin
	bstore, err := bytengine.NewByteStore("diskv", BSTORE_CONFIG)
	assert.Nil(t, err, "bst not created")
	// get bfs plugin
	mfs, err := bytengine.NewFileSystem("mongodb", BFS_CONFIG, &bstore)
	assert.Nil(t, err, "bfs not created")

	// set database
	db := "db1"

	// create directories
	err = mfs.NewDir("/var", db)
	assert.Nil(t, err, "directory not created")
	err = mfs.NewDir("/var/www", db)
	assert.Nil(t, err, "directory not created")

	// create file
	err = mfs.NewFile("/var/www/index.html", db, map[string]interface{}{})
	assert.Nil(t, err, "file not created")

	// update file
	data := map[string]interface{}{
		"title": "welcome",
		"body":  "Hello world!",
	}
	err = mfs.UpdateJson("/var/www/index.html", db, data)
	assert.Nil(t, err, "file update failed")

	// read file
	j, err := mfs.ReadJson("/var/www/index.html", db, []string{"title", "body"})
	assert.Nil(t, err, "file read failed")
	val, ok := j.(bson.M)
	assert.True(t, ok, "couldn't cast file content to bson.M")
	assert.Equal(t, val["title"], "welcome", "incorrect file content: title")
	assert.Equal(t, val["body"], "Hello world!", "incorrect file content: body")

	// copy file
	err = mfs.Copy("/var/www/index.html", "/var/www/index_copy.html", db)
	assert.Nil(t, err, "file copy failed")

	// directory listing
	list, err := mfs.ListDir("/var/www", "", db)
	assert.Nil(t, err, "directory listing failed")
	files := list["files"]
	assert.Len(t, files, 2, "file copy failed")

	// copy directory
	err = mfs.Copy("/var/www", "/www", db)
	assert.Nil(t, err, "directory copy failed")

	// directory listing
	list, err = mfs.ListDir("/www", "", db)
	assert.Nil(t, err, "directory listing failed")
	files = list["files"]
	assert.Len(t, files, 2, "directory copy failed")

	// read copied file contents
	j, err = mfs.ReadJson("/www/index_copy.html", db, []string{"title", "body"})
	assert.Nil(t, err, "file read failed")
	val, ok = j.(bson.M)
	assert.Equal(t, val["title"], "welcome", "incorrect file content: title")
	assert.Equal(t, val["body"], "Hello world!", "incorrect file content: body")
}

func TestCounters(t *testing.T) {
	// get bst plugin
	bstore, err := bytengine.NewByteStore("diskv", BSTORE_CONFIG)
	assert.Nil(t, err, "bst not created")
	// get bfs plugin
	mfs, err := bytengine.NewFileSystem("mongodb", BFS_CONFIG, &bstore)
	assert.Nil(t, err, "bfs not created")

	// set database
	db := "db1"

	val, err := mfs.SetCounter("users", "incr", 1, db)
	assert.Nil(t, err, "counter action failed")
	assert.Equal(t, val, 1, "counter action failed")

	val, err = mfs.SetCounter("users", "decr", 1, db)
	assert.Nil(t, err, "counter action failed")
	assert.Equal(t, val, 0, "counter action failed")

	val, err = mfs.SetCounter("users", "reset", 5, db)
	assert.Nil(t, err, "counter action failed")
	assert.Equal(t, val, 5, "counter action failed")

	val, err = mfs.SetCounter("user1.likes", "incr", 1, db)
	assert.Nil(t, err, "counter action failed")
	val, err = mfs.SetCounter("car.users", "incr", 1, db)
	assert.Nil(t, err, "counter action failed")

	list, err := mfs.ListCounter("", db)
	assert.Nil(t, err, "counter action failed")
	assert.Len(t, list, 3, "counter list failed")

	list, err = mfs.ListCounter("^user", db)
	assert.Nil(t, err, "counter action failed")
	assert.Len(t, list, 2, "counter list failed")
}

func TestSearch(t *testing.T) {
	// get bst plugin
	bstore, err := bytengine.NewByteStore("diskv", BSTORE_CONFIG)
	assert.Nil(t, err, "bst not created")
	// get bfs plugin
	mfs, err := bytengine.NewFileSystem("mongodb", BFS_CONFIG, &bstore)
	assert.Nil(t, err, "bfs not created")

	// set database
	db := "db1"

	// create dir and add files
	err = mfs.NewDir("/users", db)
	assert.Nil(t, err, "directory not created")
	err = mfs.NewFile("/users/u1", db, map[string]interface{}{
		"name":    "john",
		"age":     34,
		"country": "ghana",
	})
	assert.Nil(t, err, "file not created")
	err = mfs.NewFile("/users/u2", db, map[string]interface{}{
		"name":    "jason",
		"age":     18,
		"country": "ghana",
	})
	assert.Nil(t, err, "file not created")
	err = mfs.NewFile("/users/u3", db, map[string]interface{}{
		"name": "juliette",
		"age":  18,
	})
	assert.Nil(t, err, "file not created")
	err = mfs.NewFile("/users/u4", db, map[string]interface{}{
		"name":    "michelle",
		"age":     21,
		"country": "uk",
	})
	assert.Nil(t, err, "file not created")
	err = mfs.NewFile("/users/u5", db, map[string]interface{}{
		"name":    "dennis",
		"age":     22,
		"country": "france",
	})
	assert.Nil(t, err, "file not created")

	// create parser
	parser, err := bytengine.NewParser("base", "")
	assert.Nil(t, err, "parser not created")

	// search users by country
	script := `@test.select "name" "age" in /users where "country" in ["ghana"]`
	cmd, err := parser.Parse(script)
	assert.Nil(t, err, "couldn't parse script")
	rep, err := mfs.BQLSearch(db, cmd[0].Args)
	assert.Nil(t, err, "search failed")
	val, ok := rep.([]interface{})
	assert.True(t, ok, "couldn't cast search result into []interface")
	assert.Len(t, val, 2, "search failed")

	// search users by regular expression on name
	script = `
    @test.select "name" "age" in /users
    where regex("name","i") == "^j\\w*n$"`
	cmd, err = parser.Parse(script)
	assert.Nil(t, err, "couldn't parse script")
	rep, err = mfs.BQLSearch(db, cmd[0].Args)
	assert.Nil(t, err, "search failed")
	val, ok = rep.([]interface{})
	assert.True(t, ok, "couldn't cast search result into []interface")
	assert.Len(t, val, 2, "search failed")

	// search users that have a country field using 'exists'
	script = `
	    @test.select "name" "age" in /users
	    where exists("country") == true`
	cmd, err = parser.Parse(script)
	assert.Nil(t, err, "couldn't parse script")
	rep, err = mfs.BQLSearch(db, cmd[0].Args)
	assert.Nil(t, err, "search failed")
	val, ok = rep.([]interface{})
	assert.True(t, ok, "couldn't cast search result into []interface")
	assert.Len(t, val, 4, "search failed")

	// search users and return number using 'count'
	script = `@test.select "name" "age" in /users count`
	cmd, err = parser.Parse(script)
	assert.Nil(t, err, "couldn't parse script")
	rep, err = mfs.BQLSearch(db, cmd[0].Args)
	assert.Nil(t, err, "search failed")
	val2, ok := rep.(int)
	assert.True(t, ok, "couldn't cast search result into int")
	assert.Equal(t, val2, 5, "search failed")
}

func TestSetUnset(t *testing.T) {
	// get bst plugin
	bstore, err := bytengine.NewByteStore("diskv", BSTORE_CONFIG)
	assert.Nil(t, err, "bst not created")
	// get bfs plugin
	mfs, err := bytengine.NewFileSystem("mongodb", BFS_CONFIG, &bstore)
	assert.Nil(t, err, "bfs not created")

	// set database
	db := "db1"

	// create parser
	parser, err := bytengine.NewParser("base", "")
	assert.Nil(t, err, "parser not created")

	script := `
    @test.set "country"={"name":"ghana","major_cities":["kumasi","accra"]}
    in /users
    where "country" == "ghana"
    `
	cmd, err := parser.Parse(script)
	assert.Nil(t, err, "couldn't parse script")
	count, err := mfs.BQLSet(db, cmd[0].Args)
	assert.Nil(t, err, "set data failed")
	assert.Equal(t, count, 2, "set data failed")

	j, err := mfs.ReadJson("/users/u1", db, []string{})
	assert.Nil(t, err, "read file failed")
	data, ok := j.(bson.M)
	assert.True(t, ok, "couldn't cast file content to bson.M")
	country, ok := data["country"].(bson.M)
	assert.True(t, ok, "couldn't cast file content to bson.M")
	assert.Equal(t, country["name"], "ghana", "incorrect file content update")

	script = `
    @test.unset "country"
    in /users
    where exists("country") == true
    `
	cmd, err = parser.Parse(script)
	assert.Nil(t, err, "couldn't parse script")
	count, err = mfs.BQLUnset(db, cmd[0].Args)
	assert.Nil(t, err, "unset data failed")
	assert.Equal(t, count, 4, "unset data failed")

	script = `@test.select "name" in /users where exists("country") == false`
	cmd, err = parser.Parse(script)
	assert.Nil(t, err, "couldn't parse script")
	j, err = mfs.BQLSearch(db, cmd[0].Args)
	assert.Nil(t, err, "search failed")
	val2, ok := j.([]interface{})
	assert.True(t, ok, "couldn't cast search result into []interface")
	assert.Len(t, val2, 5, "search failed")
}

func TestAttachmentManagement(t *testing.T) {
	// get bst plugin
	bstore, err := bytengine.NewByteStore("diskv", BSTORE_CONFIG)
	assert.Nil(t, err, "bst not created")
	// get bfs plugin
	mfs, err := bytengine.NewFileSystem("mongodb", BFS_CONFIG, &bstore)
	assert.Nil(t, err, "bfs not created")

	// set database
	db := "db1"

	// create test file
	txt := "Hello from bst!"
	fpath := "/tmp/bfs_attach.txt"
	err = ioutil.WriteFile(fpath, []byte(txt), 0777)
	assert.Nil(t, err, "test file not created")

	data := map[string]interface{}{
		"title": "bfs test file",
		"type":  ".txt",
	}
	bfs_path := "/file_with_attachment"
	err = mfs.NewFile(bfs_path, db, data)
	assert.Nil(t, err, "file creation failed")

	// add to bfs
	_, err = mfs.WriteBytes(bfs_path, fpath, db)
	assert.Nil(t, err, "write bytes failed")

	// read from store
	fpath2 := "/tmp/bfs_attach_down.txt"
	f2, err := os.Create(fpath2)
	assert.Nil(t, err, "download test file not created")

	bstore_id, err := mfs.ReadBytes(bfs_path, db)
	assert.Nil(t, err, "read bytes failed")

	bstore.Read(db, bstore_id, f2)
	f2.Close()

	// check downloaded file data
	fdata, err := ioutil.ReadFile(fpath2)
	assert.Nil(t, err, "download test file couldn't be opened")
	assert.Equal(t, txt, string(fdata), "attachment file content has changed")
}
