# kpoward

![Go](https://github.com/goccy/kpoward/actions/workflows/test.yaml/badge.svg)
[![GoDoc](https://godoc.org/github.com/goccy/kpoward?status.svg)](https://pkg.go.dev/github.com/goccy/kpoward?tab=doc)

**K**ubernetes **Po**rt F**oward**ing utility library for Go

If you specify the name and port number of the pod on which the server is running, port forwarding will be performed using a free port in your local system.

# Motivation

When developing tools and tests, we may want to make a request from local to a server running on our kubernetes cluster. Normally, we would use `kubectl port-foward <pod name> <local port>:<remote port>` to bind a remote port to a local port and send a request to that port, but we may want to automate this set of tasks. Therefore, we made a library so that this work can be done in the Go lauguage.

# Install

```console
go get github.com/goccy/kpoward
```

# Synopsis

Create a `Kpoward` instance with [`rest.Config`](https://pkg.go.dev/k8s.io/client-go/rest#Config) and the name and port number of the pod. If necessary, use the Setter method to change a value such as `Namespace` and then call `Run` will call back the local bound free port. You can send any request to this port in the callback. Upon exiting the callback, port forwarding will automatically end and the port will be released.

```go
var (
    restCfg *rest.Config
    targetPodName = "pod-xxx-yyy"
    targetPort = 8080
)
kpow := kpoward.New(restCfg, targetPodName, targetPort)
if err := kpow.Run(context.Background(), func(ctx context.Context, localPort uint16) error {
  log.Printf("localPort: %d", localPort)
  resp, err := http.Get(fmt.Sprintf("http://localhost:%d", localPort))
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  return nil
}); err != nil {
  panic(err)
}
```

# License

MIT
