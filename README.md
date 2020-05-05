![realworld-logo](https://user-images.githubusercontent.com/13511520/81056310-4bf24b00-8f05-11ea-91d5-c98e1d6d621e.png)
---

[![Test Status](https://github.com/raahii/golang-grpc-realworld-example/workflows/test/badge.svg)](https://github.com/raahii/golang-grpc-realworld-example/actions?query=workflow%3Atest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/raahii/golang-grpc-realworld-example/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/raahii/golang-grpc-realworld-example?status.svg)](https://godoc.org/github.com/raahii/golang-grpc-realworld-example)


> ### Go/GRPC codebase containing RealWorld examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

This codebase was created to demonstrate a fully fledged fullstack application built with golang/grpc including CRUD operations, authentication, routing, pagination, and more.



## Getting started

The app listens and serves on `0.0.0.0:8080`. 


- docker-compose

  ```
  $ docker-compose up -d
  ```

  

- locally

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



- locally

  ```
  $ go test ./...
  ```

  
