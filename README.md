drainchecker is a small Go library that allows checking whether HTTP response
bodies have been fully read. Otherwise, HTTP connections cannot be reused.

USAGE
=====

Wrap the default `http.DefaultRoundTripper` and it is used by `net.Get()`,
`net.Post()` etc. You can also give it to a library:

```go
http.DefaultTransport = drainchecker.RoundTripper(http.DefaultTransport)
http.DefaultClient.Transport = http.DefaultTransport
```
