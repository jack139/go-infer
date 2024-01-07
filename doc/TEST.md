## Local test



### Start the dispatcher

```
go test -v -run TestServer
```



### Start the HTTP server for API

```
go test -v -run TestHttp
```



### Test demo API

```
cd examples
python3 test_api.py 127.0.0.1 echo
```