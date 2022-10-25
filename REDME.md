## Tools
* go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
* go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
* apt install -y protobuf-compiler

## protoc

Compiler that allows convert a protofile into a golang struct/code

## Convert protofile to golang package:

```bash
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/student.proto
```


## RPC

Stalbish the connection between the client and server and hiddens the implementation of the server.

## gRPC
It uses 
* HTTP2 -> Transport layer
* ProtoBuffers -> Data(messages) Serialize | Deserialize
* Streaming -> Allows to send data constantly through a channel