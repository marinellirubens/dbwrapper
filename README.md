# dbwrapper
## Summary
* [Project Description](#project-description)
* [Project Goals](#when-the-project-will-be-concluded)
* [References](#references)


## Project description

Project to practice golang, the idea is to create an wrapper for database requests

My goal is to create an api to do the basic operations on the database, to be inserted on legacy systems, that can do - http requests but change to include a specific database may be complicated or ad some extensive change.

At first I want to create conections with postgresql, oracle, mysql and mongodb.

Not sure yet how to structure the usage, but will try the connections first and organizing the results, returning always a list with json objects inside.

## When the project will be concluded
This section is just so a I can have the completion feeling and a goal to achieve

- When I can do on databases (postgres, oracle, mysql, mongodb, redis) the "equivalent" operations for (select, update, delete, insert)
    - [x] postgres
    - [ ] oracle
    - [ ] mysql
    - [ ] mongodb
    - [ ] redis
- [ ] JWT authentication
- [ ] log rotation implemented
- [ ] cli implemented

## references

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
