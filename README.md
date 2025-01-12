# dbwrapper
My goal is to create an api to do the basic operations on the database, to be inserted on legacy systems, that can do http requests, but change to include a specific database may be complicated.

At first I want to create conections with postgresql, oracle, mysql.

The idea is to return the result of the queries as a json, and other operations such as delete, insert, update, return only the status.

## Summary
* [Project Goals](#when-the-project-will-be-concluded)
* [Build project](#build-project)
* [Configuration](#configuration)
* [References](#references)

<!--TODO: improve this readme file-->

## When the project will be concluded
This section is just so a I can have the completion feeling and a goal to achieve

- When I can do on databases (postgres, oracle, mysql, mongodb, redis) the "equivalent" operations for (select, update, delete, insert)
    - [x] postgres
    - [x] oracle
    - [x] mysql
- [x] key authentication
- [ ] log rotation implemented (decided to use logrotate for that)
- [x] cli implemented

## Build project

To build this project you can use the following command to use the configuration on the [Makefile](Makefile) to build it
```shell
$ make build
```

This will compile an executable on the bin folder, inside the folder related to the version informed on [VERSION](VERSION)
```shell
$ cat VERSION
1.0.0
$ ls ./bin
1.0.0
$ ls ./bin/1.0.0
dbwrapper
```

## Configuration

This service uses a json file for the configuration related to the database connections and also for the server port/ip and logging file

There is an example for the configuration on the file [config.example.json](internal/config/examples/config.example.json)
```json
{
  "Server": {
    "server_port": 8080,
    "server_address": "",
    "logger_file": "/tmp/server.log"
  },
  "Databases": [
    {
      "dbid": "localdb",
      "host": "localhost",
      "port": 5432,
      "user": "myusername",
      "password": "mypassword",
      "dbname": "myusername",
      "dbtype": "postgres"
    }
  ]
}
```
You can use this file on the same folder you are executing the service, or specify the path using the parameter -f
```shell
$ ./bin/1.0.0/dbwrapper -f /opt/config.json
```

You can also check the cli options using -h
```shell
$ ./bin/1.0.0/dbwrapper -h
NAME:
   dbwrapper - A new cli application

USAGE:
   dbwrapper [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config_file value, -f value  Path for the configuration file (default: "./config.json")
   --version, -v                  Path for the configuration file (default: false)
   --help, -h                     show help
```

## References

### golang connections pooling

- https://koho.dev/understanding-go-and-databases-at-scale-connection-pooling-f301e56fa73


### nginx workers

- https://www.educba.com/nginx-worker_connections/


### deployment golang + nginx

- https://hackersandslackers.com/deploy-golang-app-nginx/
- https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-using-nginx-on-ubuntu-18-04


### creating rest api with only native modules

- https://dev.to/mauriciolinhares/building-and-distributing-a-command-line-tool-in-golang-go0
- https://tutorialedge.net/golang/creating-restful-api-with-golang/
- https://www.section.io/engineering-education/how-to-build-a-rest-api-with-golang-using-native-modules/
- https://dev.to/bmf_san/introduction-to-golang-http-router-made-with-nethttp-3nmb

### handle routing without using a framework

- https://benhoyt.com/writings/go-routing/#regex-table


### error handling good practices

- https://www.youtube.com/watch?v=7g-kGONT8ds&ab_channel=GolangDojo


### basic logging in golang:

- https://www.youtube.com/watch?v=yF7k6PxtRU8


### Anthony GG series on creating a golang REST API with almost no external packages:

- https://www.youtube.com/watch?v=pwZuNmAzaH8&list=PL0xRBLFXXsP6nudFDqMXzrvQCZrxSOm-2
