# send

This example acts like a simple `Send` method client implementation --- one message to one or more recipients over the JSON API v2 (`POST /message`). It returns the request ID, the message IDs of the accepted recipients, and the recipients rejected by the API.

## example output

```
request ID: 66666666-6666-6666-6666-666666666666
accepted recipient 0: message ID pppppppp-qqqq-rrrr-ssss-tttttttttttt
```

```
2026/08/03 10:05:45 sms manager api: 400: User is not authorized to access this resource with an explicit deny
```
