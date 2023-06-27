# requestList

Simple example on how to fetch SMS request list. List item consists of a request number (ID), gateway type/tariff used, SMS sending time, SMS send expiration time, sender name set, number of remaining receivers, request process state/status/code. 

```
go run main.go
```

## example output

```
2023/06/27 03:29:28 method: RequestList
2023/06/27 03:29:28 status: 200 OK
request ID                 : 58020767
gateway type               : high
SMS sending time           : 2023-06-26 20:45:00
SMS send exp. time         : 
sender name set            : 
remaining receivers count  : 0
request process code/status: 2 (sent)

request ID                 : 58019737
gateway type               : high
SMS sending time           : 2023-06-26 20:22:50
SMS send exp. time         : 
sender name set            : 
remaining receivers count  : 0
request process code/status: 2 (sent)

request ID                 : 58018592
gateway type               : high
SMS sending time           : 2023-06-26 20:08:36
SMS send exp. time         : 
sender name set            : 
remaining receivers count  : 0
request process code/status: 2 (sent)

```
