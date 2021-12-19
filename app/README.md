# container-api app

## Build the container image

Build image and push to local registry

```
docker build . -t canivete-api:latest
docker tag canivete-api:latest 172.19.255.222:5000/canivete-api:latest
docker image push 172.19.255.222:5000/canivete-api:latest
```

Run container

```
docker run -p 8080:8080 canivete-api:latest
```

Test container

```
http localhost:8080
```

## Run locally


Run the gin tonic server

```
go run main.go
```

## Run tests

Execute tests

```
go test ./... -cover -coverprofile=coverage.out
```

Check coverage report

```
go tool cover -html=coverage.out
```

