# Feed API challenge

# How to start DB?
- This app uses SQLite3. You don't need to install anything.

# How to build?
- Run `make build` to compile the project directly on your OS (tested on MAC M1), binary will be added to build directory, it was tested with go 1.19+.

# How to run?
- Run `make run` it was tested with go 1.19+.

# How to run tests?
- Run `make lint` to run linter.
- Run `make test` to execute unit tests.
- Run `make integration-test` to execute integration tests.

# Implementation details
- This app uses:
  - [Go 1.19](https://golang.org/doc/devel/release.html)
  - [Gin] to handle requests, responses, validations (validator V10).
  - [SQLite3] to store data.
  - [Migrate] to manage database migrations.
  - [sqlMock & Moq] to mock dependencies.
  - [Golangci-lint] to lint code.
  - [Venom] to run integration tests (WIP).
  - Valid SubReddits are stored in a txt file that is loaded at startup.

# Postman
[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/8c387cfb0e474cbdf581?action=collection%2Fimport)

# Endpoints
## CreatePost
### Description:
Method: POST
>```
>http://localhost:8082/v1/post
>```
### Body (**raw**)

```json
{
  "title": "titulo",
  "subreddit": "/r/nvidia",
  "link": "https://www.twi2tch.tv/elmariana",
  "content": "content is mututally exclusive with link, choose wisely",
  "score": 1232132,
  "nsfw": true,
  "promoted": false
}
```

### Response (**raw**)

```json
{
  "id": "35ce6d85-2260-484d-b656-24f90408efe9"
}
```

## Update Post Score
### Description:
Method: POST
>```
>http://localhost:8082/v1/post/score
>```
### Body (**raw**)

```json
{
  "id":"2bf5b210-b4f4-47b3-bc7e-bd5034bf129d",
  "score":45
}
```

### Response (**raw**)

```json
{
  "error": null
}
```

## Get Feed
### Description:
Method: GET
>```
>http://localhost:8082/v1/feed?page=1&limit=10
>```

### Response (**raw**)

```json
{
  "data": [
    {
      "id": "88539220-1d68-42de-aa5c-150115ddcecd",
      "title": "titulo",
      "author": "t2_Ej6sNFr4",
      "subreddit": "/r/nvidia",
      "link": "https://www.twi2tch.tv/elmariana",
      "score": 1232132,
      "promoted": false,
      "nsfw": true
    },
    {
      "id": "35ce6d85-2260-484d-b656-24f90408efe9",
      "title": "titulo",
      "author": "t2_rcXbzrcA",
      "subreddit": "/r/nvidia",
      "link": "https://www.twi2tch.tv/elmariana",
      "score": 1232132,
      "promoted": false,
      "nsfw": true
    }
  ],
  "pagination": {
    "pageCount": 2,
    "next": "/v1/feed?page=2",
    "page": 1,
    "limit": 2,
    "total": 62
  }
}
```


_________________________________________________
Author: [Carlos Flores](https://github.com/carfloresf)