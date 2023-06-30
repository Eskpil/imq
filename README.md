# IMQ (Integrated message queue)

IMQ is a library which allows developers to embed a message queuing
system in their applications without having dependencies on other
technologies like Kafka or RabbitMQ. 

IMQ is technically just a library which you include and start up inside
your api, or a potential microservice. It travels over TCP/TLS allow it
to be reverse proxied through services like Nginx. Although, you may
need a custom Nginx module.

```go
// insert example server code here
```

The server operates with the same API as a client, but instead of
connecting over a socket, it operates through the broker system. This
allows to server both to listen, publish and call queues.

```go
// insert example client code here
```
