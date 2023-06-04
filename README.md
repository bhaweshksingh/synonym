# Synonym service

## Pre requisites

- Docker
- Golang v1.14+

## Running App

1. Build and run the app container.

`make app`

5. Inspect logs using docker

`docker logs synonym-service-go -f`

### Setup:

Run the queries:

## Verifying the Functionality

### Add synonyms

```shell
curl --location 'localhost:8888/synonyms' \
--header 'Content-Type: application/json' \
--data '{
    "synonyms": ["abc","ghi"]
}'
```

```shell
curl --location 'localhost:8888/synonyms' \
--header 'Content-Type: application/json' \
--data '{
    "synonyms": ["def","ghi"]
}'
```

### Search Synonyms

```shell
curl --location 'localhost:8888/synonyms/search?word=abc'
```
