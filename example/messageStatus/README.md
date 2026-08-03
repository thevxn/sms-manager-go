# messageStatus

The full state of a single sent message, looked up by its `message_id` over the REST API v1 (`GET /message`). Replaces the former `RequestStatus` method.

## example output

```
phone number   : 420123456789
send status    : sent
delivery status: delivered
delivered at   : 2026-07-30T10:00:04Z
price          : 1.20
```
