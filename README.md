# A Go framework of DL model inference and API delployment



## 测试



### 启动 dispatcher

```
go test -v -run TestServer
```



### 启动 http

```
go test -v -run TestHttp
```



### 测试脚本

```
cd examples
python3 test_api.py 127.0.0.1 _
```