# pm-grpc

PM GRPC is a proxy which allows HTTP API tools like Postman to interact with gRPC servers.

## Requirements
- grpc+json codec must be enabled on the grpc server
- Postman must be configured to use the proxy

Configuration of the proxy and its dependencies is a three step process.

Register a JSON codec with the gRPC server. In Go, it can be automatically registered simple by adding the following import:

`import _"github.com/jnewmano/grpc-json-proxy-codec"`

2. Run the grpc-json-proxy.

`go get -u github.com/jnewmano/grpc-json-proxy`
`grpc-json-proxy`

3. Configure Postman to send requests through the proxy.
Postman -> Preferences -> Proxy -> Global Proxy

`Proxy Server: localhost 7001`

![Postman Proxy Configuration](https://cdn-images-1.medium.com/max/1600/1*oc09cwpCC9XrjpU9Gl5YTw.png)

Inspired by Johan Brandhorst's [grpc-json](https://jbrandhorst.com/post/grpc-json/)
