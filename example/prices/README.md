# prices

The price list by channel and destination (`GET /credit/prices`). Replaces the former `GetPrice` method --- the v2 APIs quote destination prices, not the price of a message text that has not been sent yet.

## example output

```
destination MCC: 230
price per SMS  : 0.89 CZK / 0.040 EUR
operator 1     : 0.89 CZK / 0.040 EUR
```
