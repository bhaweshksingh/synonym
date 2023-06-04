# Synonym service

The Thesaurus Service is designed to handle and manage sets of synonyms.
It allows users to add synonyms and search for synonyms of a word.
It uses a thread-safe design that supports concurrent requests to add and search synonyms.

## Pre requisites

- Docker
- Golang v1.14+

## Running App

1. Build and run the app container.

`make app`

5. Inspect logs using docker

`docker logs synonym-service-go -f`

## Design:

The Thesaurus Service uses the concept of an adjacency list to store the synonyms of each word and a map to store the
word's corresponding graph ID. The adjacency list is a map where the key is a word and the value is another map storing
the synonyms of the word as keys (values are not used in this case).

The graph ID map stores the word as the key and the corresponding graph ID as the value. If words are synonyms, they
will share the same graph ID.

The adjacency list and the graph ID map are protected by a Mutex to prevent data races.

The Thesaurus Service also uses buffered channels to handle requests to add synonyms and search synonyms, as well as to
return search results.

There is a separate goroutine running the main logic of the service. It listens on the channels for incoming requests,
processes them, and sends the results back through the result channel.

### Capabilities

- Add new sets of synonyms: For example, adding a pair of synonyms such as “begin” and “start”. You can also add more
  than
  two synonyms at a time.

- Search for synonyms: In the above example, searching for either “begin” or “start” will return the respective synonym.
  All synonyms of a word will be returned.

- Concurrent requests: The service is designed to be thread-safe and can handle concurrent requests.

- Transitive rule: For example, if “A” is added as a synonym for “B", and "B" is added as a synonym for “C", then
  searching the word “C" should return both “B" and "A".

### Limitations

The service uses simple in-memory data structures, so the data is not persistent. When the service restarts, all
previously added synonyms will be lost. Also, because the service uses a goroutine to handle requests, if a request is
sent but the service restarts before the request is processed, that request will also be lost.

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

