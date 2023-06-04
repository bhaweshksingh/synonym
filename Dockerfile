FROM golang:alpine as builder
WORKDIR /synonym
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o synonym cmd/*.go

FROM scratch
COPY --from=builder /synonym/synonym .
ENTRYPOINT ["./synonym","http-serve"]
