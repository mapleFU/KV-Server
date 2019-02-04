# KV-Server
A KV Server.

## Quickstart

First, you should clone this project in your directory. 

```bash
git clone https://github.com/mapleFU/KV-Server
cd KV-Server/server
go run main.go
```

You should first run `main.go` in `server` package. And in client package, you are allowed to use grpc to connect and send request to the server.

## Targets

### Server

* [x] Get
* [x] Set
* [x] Delete
* [x] Scan
* [x] Client/Server tools
* [x] Persistence/Recover
* [ ] Compaction
* [ ] Benchmark

### Client

Implemented interface like server.

## Design

### Protocol

We use grpc to design our sevice. In `proto/kv-server.proto`, we design `Get` `Set` `Del` and `Scan`.  And client can use this to intereact with server.

```protobuf
service KVServicer {
    rpc Get(Key) returns (Value) {}
    rpc Set(KVPair) returns (OperationOK) {}
    rpc Delete(Key) returns (OperationOK) {}
    rpc Scan(ScanArgs) returns (Values) {}
}
```

## storage

We use [bitcask](http://basho.com/wp-content/uploads/2015/05/bitcask-intro.pdf) as our storage model, it was used in basho and beandb. 





