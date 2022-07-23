# Go实现antigen推理和api服务



## 测试



### 编译

```
make
```



### 启动 dispatcher

```
build/antigen-go server 0
```



### 启动 http

```
build/antigen-go http
```



### 测试脚本

```
python3 test_api.py 127.0.0.1 _
```



### 压力测试
```
python3 stress_test.py 1 1
```
