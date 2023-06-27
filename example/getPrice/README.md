# getPrice

Simple example on how to fetch the information about SMS costs for given phone number(s), and given message length. Returns number of valid receivers (phone numbers), number of SMS messages to be sent (text could be split into more messages), message's character count, price for one SMS messages, and the price for all possible SMS messages to be sent.

```
go run main.go
```

## example output

```
2023/06/27 09:38:42 method: GetPrice
2023/06/27 09:38:43 status: 200 OK
valid receivers count    : 1
count of SMS msgs        : 4
SMS msg's character count: 574
price per one SMS msg    : 0.75 CZK
price per all SMS msgs   : 3.00 CZK
```
