# getUserInfo

Simple example on how to fetch and print SMS Manager user's information (credit balance, used gateway type, and used sender name).

```
go run main.go
```

## example output

```
2023/06/27 03:36:20 method: GetUserInfo
2023/06/27 03:36:20 status: 200 OK
credit left : 4.00 CZK (exc. VAT)
sender name : SmsManager
gateway type: high
```

```
2023/06/27 09:39:52 method: GetPrice
2023/06/27 09:39:52 status: 401 Unauthorized
response error code: 103
response error msg: The username or password you entered is incorrect
```
