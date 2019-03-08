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
* [x] Compaction
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

We use [bitcask](http://basho.com/wp-content/uploads/2015/05/bitcask-intro.pdf) as our storage backend, it was used in basho and beandb. And this is interface for our backend:

```go
// Engine is an interface for storage backend like LSM-tree or Bitcask
//
type Engine interface {
	// Get the data in
	Get([]byte) ([]byte, error)
	// Scan the data of engine, return bytes array
	Scan(cursor ScanCursor)	([][]byte, error)
	// Put the data into the engine
	Put([]byte, []byte) error
	// delete the data in the engine
	Del([]byte) error
}
```

In this Key/Value server, I only implemented bitcask.

## Bitcask

Bitcask is our storage backend. The logic of bitcask is like the following image:

![屏幕快照 2019-03-07 下午2.41.24](https://nmsl.maplewish.cn/blog:/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/README.md:屏幕快照 2019-03-07 下午2.41.24.png)

It has only one active data file and many older data files. The data will be store in all data files, but the write or delete operations will be only append to active data file. The data model of data file is:

 ![屏幕快照 2019-03-07 下午2.45.56](https://nmsl.maplewish.cn/blog:/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/README.md:屏幕快照 2019-03-07 下午2.45.56.png)

In my implement, you can see:

```go
 /**
 header
| crc | tstamp | ksz | value_sz |
  */
type BitcaskStoreHeader struct {
	Crc uint32
	TimeStamp uint32
	//timeStamp uint64	// u64 unixNano
	KeySz uint32
	ValueSz uint32
}
```

To find the data in datafile, we keep a `keydir` in memory. When we want to read value from datafile with the key, we can first lookup the `file_id` and the file position in keydir.

![屏幕快照 2019-03-07 下午2.49.42](https://nmsl.maplewish.cn/blog:/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/README.md:屏幕快照 2019-03-07 下午2.49.42.png)

So, our implement can be design as follows:

* Get:
  1. read file_id and bios from keydir
  2. get the data in the file
  3. Unmarshall the bytes and get the value
* Put
  1. append data to active data file
  2. append the record to keydir, and after this we can write aof log.
  3. Put success.
* Delete
  1. append an record with value_sz is 0 to active data file
  2. delete the record in keydir, write aof log.
  3. Delete success.

Specifically, our scan works just like redis, we have `ScanCursor`:

```go
type ScanCursor struct {
	// cursor
	Cursor int
	// match
	UseMatchKey bool
	MatchKeyString string
	// count has a default: 10
	Count int
}
```

After scan , the cursor and the values will be return.

### Usage

As for the usage of bitcask.

```go
bitcask := Open(workDir, &options.Options{UseLog:true})
```

You can create a bitcask object like this. You may notice `options`, it's our config for bitcask. It contains `MaxFileSize` for data file. If you want to use the default config, just use `Open(workDir, nil)`.

Now, you can use crud for the bitcask k-v storage engine.

### aof log and hint file

If our program crash, we will lose our keydir. So we need to recover after restart/crash. In our options, if you use `UseLog = true`, we will use aof log. If write success, data will be append to logfile and execute `fsync`. And if we close the program, we will generate a hint file, which was a backup of keydir.

### switch file

In our program, if the current data file size is larger than `MaxFileSize` (default is 1024 * 500, you can adjust it), the we will change the data file, and create a new datafile.

### sync mode

The following sync strategies are available:

- `none` — lets the operating system manage syncing writes (default)
- `o_sync` — uses the `O_SYNC` flag, which forces syncs on every write
- Time interval — Riak will force Bitcask to sync at specified intervals

In `BitcaskBufferManager`, we will set different writer for them. 

### compaction

In bitcask, if too many datafile will generated, read from them will me a mess. So we can use compaction to merge old datafiles.

