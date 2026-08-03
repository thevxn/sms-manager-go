# flags

A simple example on how to use implemented flags package support. This way it is very easy to use a built binary in various scripts, workflows, cronjobs etc.

Two methods are implemented --- `send` sends a single message, `credit` prints the credit balance, or reports it to a Prometheus PushGateway when `-pushGatewayURL` is given.

```bash
go run main.go -method send -apiKey xxx -phoneNumber 420123456789 -body "hello world!" -sender MyCompany -gateway high
```

```
request ID: 66666666-6666-6666-6666-666666666666
message ID: pppppppp-qqqq-rrrr-ssss-tttttttttttt
```

```bash
go run main.go -method credit -apiKey xxx
```

```
150.50
```
