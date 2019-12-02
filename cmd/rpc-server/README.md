Run an `spt` render server:

```go
go run main.go
2019/12/03 00:01:51 rpc-server ready
```

Connect to it like this:

```go
Render("test.png", testScene(), []Renderer{
	NewRPCRenderer("<ip>:<port>"),
})
```

Example taken from [spt_test.go](https://github.com/seanpringle/spt/blob/master/spt_test.go).

Now start up a whole cluster of them!