![realworld-logo](https://user-images.githubusercontent.com/13511520/81056310-4bf24b00-8f05-11ea-91d5-c98e1d6d621e.png)
---

[![Test Status](https://github.com/raahii/golang-grpc-realworld-example/workflows/test/badge.svg)](https://github.com/raahii/golang-grpc-realworld-example/actions?query=workflow%3Atest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/raahii/golang-grpc-realworld-example/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/raahii/golang-grpc-realworld-example?status.svg)](https://godoc.org/github.com/raahii/golang-grpc-realworld-example)


> ### Go/GRPC codebase containing RealWorld examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)

This codebase was created to demonstrate a fully fledged fullstack application built with golang/grpc including CRUD operations, authentication, routing, pagination, and more.



## How it works

![figure](https://user-images.githubusercontent.com/13511520/81163457-bff62700-8fc9-11ea-897a-60f27d9f9c8b.png)


- Using [grpc-gateway](https://grpc-ecosystem.github.io/grpc-gateway/) as reverse-proxy server which translates a **RESTful JSON API into gRPC**.

- Using **Go** to implement realworld backend server.

  - grpc: [grpc-go](https://github.com/grpc/grpc-go)
  - auth token: [jwt-go](https://github.com/dgrijalva/jwt-go)
  - ORM: [gorm](https://github.com/jinzhu/gorm)
  - logging: [zerolog](https://github.com/rs/zerolog)

- Using **MySQL** to store data.

  

## Getting started

The app listens and serves on `0.0.0.0:8080`. 


- docker-compose

  ```
  $ docker-compose up -d
  ```

  

- local

  - Install Go 1.13+, MySQL
  - set environment variables to connect database [like this](https://github.com/raahii/golang-grpc-realworld-example/blob/master/env/local.env).

  ```
  $ go run server.go # run grpc server
  $ go run gateway/gateway.go # run grpc-gateway server
  ```



## Test

- docker-compose

  ```
  $ docker-compose run app go test ./...
  ```



- local

  ```
  $ go test ./...
  ```




## TODOs

- [x] Users and Authentication
  - [x] `POST /user/login`: Existing user login
  - [x] `POST /users`: Register a new user
  - [x] `GET /users`: Get current user
  - [x] `PUT /users`: Update current user

- [x] Profiles
  - [x] `GET /profiles/{username}`: Get a profile
  - [x] `POST /profiles/{username}/follow`: Follow a user
  - [x] `DELETE /profiles/{username}/follow`: Unfollow a user

- [ ] Articles
  - [ ] 