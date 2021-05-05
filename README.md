# Robotoscope

A silly web app used for some dokku experiments.

When deployed it serves a /robots.txt and remembers which user agents
request it. Getting "/" serves a report listing which user agents have
been encountered.


## Local Testing

```shell
$ docker run -d --name robotoscope-db -p5432:5432 -e POSTGRES_PASSWORD=hunter2 postgres
$ export DATABASE_URL=postgres://postgres:hunter2@localhost:5432/postgres
$ go build
$ ./robotoscope &
[1] 5601
$ curl -s http://localhost:5000/robots.txt
User-agent: *
Disallow: /secret/
$ curl -s http://localhost:5000/list.txt
  1: "curl/7.64.1"
$ kill %1
[1]  + 5601 terminated  ./robotoscope
$ docker stop robotoscope-db
```
