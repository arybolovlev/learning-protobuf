# learning-protobuf

## Regenerate gRPC code

```console
$ cd proto
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative number.proto
$ cd -
```

## Run

### Run server

```console
$ go run main.go
```

### Run client

```console
$ cd client/
$ go run client.go
```
