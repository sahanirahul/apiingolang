## apiingolang
API in Go

This Go project consists of 2 component:
1. API server (serving a *GET API* that returns 3 unique activities every time it is called as response)
2. DB worker (this inserts the fetched activities in postgres database)

#### postgresql queires
``` sql
- create database test

- create table activity (
    id serial primary key,
    activity_key varchar(16) not null,
    activity_content VARCHAR(512),
    created_at timestamp default current_timestamp
  );
```
#### api end point
- `/api/public/vi/activities` 
- the api serves on `port=9000`

#### how to run
- As Docker Container
  - Run the start.sh (./start.sh) script file to start, and stop.sh (./stop.sh) file to stop
- On local machine
  - Run `go run .`
 *Postgres DB connection details should be added in config/config.local.json file*


#### worker pool
- the program contains a generic worker pool implementation in go for parallel processing of request. At the same time the workers can be configured as per requirement. This is better than simply spawning go routines, if we were to spawn multiple go routine for a single request, the system can crash when there is a outburst of request. Also spawning a new go routine is costlier than a worker go routine that processes task by fetching from a job queue (channel)
- Here we have 2 worker pool:
  - First with 3 worker that is used only for calling [https://www.boredapi.com/api/activity] to fetch activities, this restricts out application to hit at max 3 times at any given point.
  - Second is a general worker pool used to execute task like db insert.
