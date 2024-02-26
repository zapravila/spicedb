```
./spicedb --log-format=auto --log-level=trace serve \
--grpc-preshared-key=123456789 --datastore-engine=postgres \
--datastore-conn-uri=postgres://solomn_user:solomonpass@127.0.0.1:30432/rulesdb
```