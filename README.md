# dbwrapper
Project to practice golang, the idea is to create and wrapper for databases requests

My goal is to create an api to do the basic operations on the database, to be inserted on legacy systems, that can do - http requests but change to include a specific database may be complicated or ad some extensive change.

At first I want to create conections with postgresql, oracle, mysql and mondodb.

Not sure yet how to structure the usage, but will try the connections first and organizing the results, returning always a list with json objects inside.

## references

golang connections pooling

- https://koho.dev/understanding-go-and-databases-at-scale-connection-pooling-f301e56fa73


nginx workers

- https://www.educba.com/nginx-worker_connections/


deployment golang + nginx

- https://hackersandslackers.com/deploy-golang-app-nginx/
- https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-using-nginx-on-ubuntu-18-04


creating rest api with only native modules

- https://dev.to/mauriciolinhares/building-and-distributing-a-command-line-tool-in-golang-go0
- https://tutorialedge.net/golang/creating-restful-api-with-golang/
- https://www.section.io/engineering-education/how-to-build-a-rest-api-with-golang-using-native-modules/
- https://dev.to/bmf_san/introduction-to-golang-http-router-made-with-nethttp-3nmb


error handling good practices

- https://www.youtube.com/watch?v=7g-kGONT8ds&ab_channel=GolangDojo


basic logging in golang:

- https://www.youtube.com/watch?v=yF7k6PxtRU8

