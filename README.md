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

## Test using perf tool Vegeta

### Install Vegeta

```shell
brew install vegeta
```

### Search REq

```
echo "GET http://localhost:8888/synonyms/search?word=word2" | vegeta attack -rate 1000 -duration 1m | vegeta report
```


### Add Req

```shell
echo "POST http://localhost:8888/synonyms" | vegeta attack -body synonyms.json -rate 1000 -duration 1m | vegeta report
```

