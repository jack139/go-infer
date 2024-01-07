## Tensorflow running environment (1.15.4)



### Install tensorflow’s C library
CPU
```
sudo tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-1.15.0.tar.gz
```
GPU
```
sudo tar -C /usr/local -xzf libtensorflow-gpu-linux-x86_64-1.15.0.tar.gz
```

```
sudo ldconfig
```

CentOS needs:
```
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
```



### Install tensorflow's go library
```
go get github.com/tensorflow/tensorflow/tensorflow/go@v1.15.4
go test github.com/tensorflow/tensorflow/tensorflow/go
```
