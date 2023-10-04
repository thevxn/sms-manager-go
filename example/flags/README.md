# flags

A simple example on how to use implemented flags package support. This way it is very easy to use a built binary in various scripts, workflows, cronjobs etc.

The example's response is trimmed to print only the credit balance for the GetUserInfo method.

```bash
go run main.go -method GetUserInfo -apiKey xxx -phoneNumber 420123456789 -gatewayType high
```

```
123.45
```
