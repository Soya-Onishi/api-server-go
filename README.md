# Backend API Server Exercise

This is a tutorial repository for api server with Go.  
This api server manage todo task list.

I got knowledge from this tutorial like below.

- Docker fundamental part(container, volume, network)
- Docker compose
- MySQL fundamental(DML, DDL, etc...)
- Go fundamental part(syntax, project structure, unit test, etc...)
- fundamental part of gin(framework for Go)
- fundamental part of cURL

## How to run

To run this project, Docker and docker-compose is required.

```
cd <repository root>/build
docker-compose up -d
```

After db container is up, you can send request to api-server.

To check whether db is up,

```
docker container log api-server-db
```

To get todos from api-server

```
curl -X GET "localhost:8080/todos"
```

You can post new todo to api-server

```
curl -X POST "localhost:8080/todos" -H "application/json" -d '{ "id": "4", "name": "new todo" }'
```

And you can also delete todo

```
curl -X DELETE "localhost:8080/todos?id=1"
```
