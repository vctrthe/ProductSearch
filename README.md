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