## postman.tmpl

This method of Postman Collection generation relies on the documentation generator for `protoc`: https://github.com/pseudomuto/protoc-gen-doc

Install with:
```
go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
```

Download the `postman.tmpl` file locally.

Building the collection for `helloworld.proto`:

```bash
POSTMAN_COLLECTION_NAME="Hello World gRPC JSON API" protoc \
--doc_out=./ \
--doc_opt=./postman.tmpl,helloworld_collection.json \
helloworld.proto
```

Try importing one of these example collections into Postman.  View the documentation and requests to see all populated fields.

Hello World Collection: https://gist.githubusercontent.com/kevinswiber/867699e4efe57ccde83fb510c9de6075/raw/b082f5505b22112d5f989dae2765a065de2e53fd/zz_helloworld_collection.json

Route Guide Collection: https://gist.githubusercontent.com/kevinswiber/867699e4efe57ccde83fb510c9de6075/raw/b082f5505b22112d5f989dae2765a065de2e53fd/zz_route_guide.proto
