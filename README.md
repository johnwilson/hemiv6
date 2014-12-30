# Hemi V6: Bytengine for Heroku

## About

**[Bytengine](https://github.com/johnwilson/bytengine "Bytengine")** is a scalable content 
repository built with Go. Its API is accessible from any Http client library so 
you can start coding in your favorite language!

Hemi V6 was created to make Bytengine easily deployable to Heroku for testing and
production use.

## Quick deployment

**Step 1** Create a Redis instance with [Redis To Go](https://redistogo.com/). Your instance
connection url will be of the form:

```
    redis://redistogo:<password>@<host>:<port>/
```

**Step 2** Create a Mongodb database with [Mongolab](https://mongolab.com/) and **give your database
name a `bfs_` prefix**. After creating a database user, your connection url would be:

```
    mongodb://<user>:<password>@<host>:<port>/<database>
```

**Step 3** Deploy to heroku [![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

Environment variables:

* *AUTH_DB* & *STORE_DB*: Mongolab `<database>` (eg. `bfs_testdb`)

* *REDIS_URL*: Redis to go url: `redis://redistogo:<password>@<host>:<port>/`

* *MONGODB_URL*: Mongolab url: `mongodb://<user>:<password>@<host>:<port>/<database>`

**Step 4** Connect to your Bytengine instance with the following default user
credentials (please change after login):

* username: `admin`
* password: `password`
* port: `80`
* host: `<your_heroku_app.herokuapp.com>`

**Step 5** Create (initialize) the database with the following bytengine command
(remove the `bfs_` prefix from database name):

```
    server.newdb "<database_name_without_prefix>"
```

## Help

If you have any problems, submit an issue or create a Stackoverflow question with 
`bytengine` tag.