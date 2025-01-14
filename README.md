This codebase was created to demonstrate a fully fledged fullstack application built with 

 - [Fiber v2](https://github.com/gofiber/fiber/tree/v2.52.6) an Express inspired web framework for Go
 - [HTMX](https://htmx.org/) to connect the frontend (html + js) with the backend
 - [Slug](https://github.com/gosimple/slug) for user friendly URLS
 - [OR Mapper gorm](gorm.io/gorm) and a [Go native driver for GORM to sqlite](https://github.com/glebarez/sqlite)
 - and [other packages](https://github.com/agmadt/realworld-fiber-htmx/blob/main/go.mod)

that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec

## Project Overview

"Conduit" is a social blogging site (i.e. a Medium.com clone). It uses a custom API for all requests, including authentication.

# Installation
```
1. clone this repository
2. run: go get
3. this project is using sqlite and its already seeded, look at conduit.sqlite
4. run: go run main.go
6. use test@email.com|secret for logging in
	6.1. or can register from the web
```
