# send

This example acts like a simple Send method client implementation. It returns the indication of request status (OK, or ERROR), request number/ID, phone number(s) used for the given request,  and the optional custom ID (if set by request).

## example output

```
2023/06/27 10:05:45 method: Send
2023/06/27 10:05:46 status: 200 OK
status         : OK
request ID     : 58059796
phone number(s): 420123456789
custom ID      : 123456
```

```
2023/06/27 09:39:52 method: Send
2023/06/27 09:39:52 status: 401 Unauthorized
response error code: 103
response error msg: The username or password you entered is incorrect
```
