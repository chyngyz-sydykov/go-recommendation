![example workflow](https://github.com/chyngyz-sydykov/go-recommendation/actions/workflows/ci.yml/badge.svg)
![Go Coverage](https://github.com/chyngyz-sydykov/go-recommendation/wiki/coverage.svg)

# About the project

This is a one of the microservices for personal pet project as study practice. the whole system consists of 3 microservices:
 - **go-web** works as an api gateway. the endpoints include CRUD actions for book, create endpoint for saving rating
 - **go-rating** another microservice that saves new rating and return list of rating by book id. the communication between go-web and go-rating is via gRPC [link](https://github.com/chyngyz-sydykov/go-rating)
- **go-recommendation** (current project) a microservice that holds business logic related with the recommendation of books depending on rating and how often the book is edited or created. the communication will be done via RabbitMQ

# Installation

 - clone the repo
 - install docker
 - copy `.env.dist` to `.env`
 - run `docker-compose up --build`

# Testing

On initial project setup, please manually create a database for tests. check the database name in env.test file. to run use following commands:

run tests `APP_ENV=test go test ./tests/`

run tests without cache `go test -count=1 ./tests/`

run tests within docker (preferred way) `docker exec -it go_rest_api bash -c "APP_ENV=test go test -count=1 ./tests"`

run test coverage on local machine `docker exec -it go_rest_api bash "scripts/coverage.sh"`
`go tool cover -html=coverage/filtered_coverage.out`

# Handy commands

To install new package

`go get package_name`

to clean up go.sum run

`go mod tidy`

to run test

running project via docker
`docker-compose up --build`
`docker-compose down`

`docker-compose logs -f`