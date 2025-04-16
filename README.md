# ProductSearch

## Project Structure
```
/ProductSearch
│
├── main.go
│
├── config/
│   │
│   └── elastic.go              -> Elasticsearch config/init
│
├── models/
│   │
│   └── product.go              -> Struct for your product
│
├── services/
│   │
│   └── search_service.go       -> Logic to search in Elasticsearch
│
├── handlers/
│   │
│   └── search_handler.go       -> HTTP Handler for search API
│
├── utils/
│   │
│   └── csv_loader.go           -> Load CSV to Elasticsearch
│
└── routes/
    │
    └── router.go               -> Route Setup
```
## Packages Used
```
gin
swagger
go-elasticsearch
```

## Run
To run the project :
```
go run main.go
```
or to run the project with hot-reload feature
```
air .
```

## API Endpoints
```
localhost:8080/search?q=            -> Main Search Endpoint
localhost:8080/docs/index.html      -> API Documentation
```

## Generating Documentation
I use `swag` package to use `swag init` command. (`go install github.com/swaggo/swag/cmd/swag@latest`). To generate, firstly I need to make comments in the `main.go` and `search.go` (the controller/handler).

Comment examples:
- main.go
```
// @title Product Search API
// @version 1.0
// @description This is a product search API using Elasticsearch
// @host localhost:8080
// @BasePath /
```

- search.go
```
// @Summary Search for products by keyword or full-text
// @Description Search for products using Elasticsearch. Results are sorted by score in ascending order. (Lowest to highest score)
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} model.SearchResponse
// @Failure 400 {object} object "Bad Request"
// @Router /search [get]
```

Then on the root project directory, I typed `swag init -g main.go` (`g` flag is to point `swag init` to read `swagger general API Info` in main.go, more info at [here](https://github.com/swaggo/swag))

## Elasticsearch Docker Compose
I've also provided a Docker compose file for deploying Elasticsearch instance on local host. Please change the environment as needed.

For `config/elastic.go`, please change the `Addresses`, `Username`, and `Password` to match the needs.
