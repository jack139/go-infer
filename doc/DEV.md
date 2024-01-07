## Development Guide

The following is a brief description of the process of using go-infer to quickly build a model inference API and deploy it. For specific codes, please refer to the complete examples in the [examples](../examples) directory.



### 1. Design ideas

The request processing flow is as shown below. The go-infer framework hides common logic such as API concurrent request processing, inference request serialization queuing, and Dispatcher distribution services. When users develop, they only need to process the parameters of the API request (API entry function) and the model inference part (model inference function ).

<img src="arch2.png" alt="Calling process" width="300" />



In the code implementation, the framework defines a simplified model interface. As long as the user implements the method functions under it, the framework can implement the process in the above figure. The model interface is defined as follows:

```go
// Model interface definition
type Model interface {
	ApiPath() (string) // HTTP URL
	ApiEntry(*map[string]interface{}) (*map[string]interface{}, error)  // Process of handling API parameters
	Init() (error)  // Model initialization, loading weights, etc.
	Infer(string, *map[string]interface{}) (*map[string]interface{}, error)  // The process of model inference
}
```



- ApiPath() is simple. It returns the URL path string of the API. When the HTTP server is started, the URL service will be registered based on this string.
- ApiEntry() is used to process the parameters passed in by the API, and usually checks the validity of the parameters based on business logic. The input parameter is a key-value map, including the content of the data field passed in by the API (for the API input parameter structure, please refer to [API Document Template](API.md)), and the return value is also a key-value map, including the content passed to inference. The content of the function.
- Init() is used to load model weights and model initialization related work. It will be called when the Dispatcher server starts.
- Infer() is used to implement specific inference services. The input parameter is the parameter data processed by requestId and ApiEntry(), and the output parameter is the content returned in the data field in the API return result.



For specific examples, please refer to [code example](../examples/models/embedding/bert_embedding.go)



### 2. Command line integration

The framework provides command line integration that can be integrated into the user's command line instructions:

```go
// Add model instance
types.ModelList = append(types.ModelList, &embedding.BertEMB{})

// Command line settings
rootCmd.AddCommand(cli.HttpCmd)
rootCmd.AddCommand(cli.ServerCmd)
```

Before adding the command line, add the model interface implemented above to the ModelList of the framework. The framework will perform initialization and registration, and find the corresponding model during request processing.



For details, please refer to [code example](../examples/main.go)



### 3. Configuration

For specific content, please refer to [configuration file example](../examples/config/settings.yaml)。

The configuration file path defaults to ```config/settings.yaml```. You can use ```--yaml``` to specify the configuration file path in the command line server and http command parameters.



### 4. Model weight export

#### (1) Tensorflow

Please refer to [export_tf_bert.py](../examples/export/export_tf_bert.py)



#### (2) Keras

Please refer to [export_keras_cnn.py](../examples/export/export_keras_cnn.py)



### 5. Deployment

#### (1) Local testing

Compiling

```bash
cd examples
make	
```



Start the Dispatcher and inference service

```bash
build/go-embedding server 0
```



Start HTTP API service

```bash
build/go-embedding http
```



API testing

```bash
python3 test_api localhost mobile
```



#### (2) Distributed architecture deployment

The go-infer framework has implemented serialized execution of API concurrent processing (Http server) and inference module (Dispatcher server). Therefore, in actual application deployment, the main consideration is the processing capability during peak concurrency periods. When concurrency is not high during peak periods, you can choose stand-alone deployment, that is, deploy the HTTP server and Dispatcher server on the same physical server. When the amount of concurrency increases during peak periods, there are three main options to improve concurrent processing capabilities (only the deployment of the CPU environment is considered here):

1. Increase the computing power of a single server.
2. Http server, Dispatcher server and redis are deployed on 3 different servers respectively.
3. Based on option 2, analyze the computing power bottleneck and horizontally expand the HTTP server and Dispatcher respectively.



Here we take option 3 as an example for demonstration. Please refer to the following figure for the deployment architecture:

<img src="arch.png" alt="Distributed deployment architecture" width="300" />



First, assume that the deployment environment is as follows:

1. One nginx server (192.168.0.100) serves as the API request entry and performs load balancing
2. One redis server (192.168.0.101) (Considering system stability, a redis cluster can be deployed. Please refer to the redis documentation for details)
3. Two HTTP servers (192.168.0.102, 192.168.0.103) provide API processing services
4. Four Dispatcher servers (192.168.0.104, 192.168.0.105, 192.168.0.106, 192.168.0.107)

Among them, in order to distribute API requests to the inference server evenly, two redis queues are established, each HTTP server is associated with one queue, and each queue backend is associated with two Dispather servers. In this way, concurrent requests of each HTTP server are processed by two Dispatcher servers. Moreover, each Dispatcher server sets the MaxWorkers parameter according to the number of CPU cores. (assuming the server is 8 cores)



> Note: The following is only a fragment of the configuration file, other content needs to be completed.



##### 1. nginx.conf

```nginx
upstream goinfer {
    least_conn;
    server 192.168.0.102:5000;
    server 192.168.0.103:5000;
}


server {
    listen 5000;
    location / {
      proxy_pass http://goinfer;
    }
}
```



##### 2. redis.conf

```redis
bind 192.168.0.101
port 7480
requirepass e18ffb7484f4d69c2acb40008471a71c
client-output-buffer-limit pubsub 32mb 8mb 60
```



##### 3. settings.yaml shared configuration

```yaml
# HTTP server parameters
API:
    Port: 5000
    Addr: 0.0.0.0
    SM2PrivateKey: "JShsBOJL0RgPAoPttEB1hgtPAvCikOl0V1oTOYL7k5U=" # SM2 private key
    AppIdSecret: { # The appid and secret assigned for the API calls
        "3EA25569454745D01219080B779F021F" : "41DF0E6AE27B5282C07EF5124642A352",
    }

# Parameters for the inference service queue
Server:
    RedisServer: "127.0.0.1:7480"
    RedisPasswd: "e18ffb7484f4d69c2acb40008471a71c"
    MessageTimeout: 10 # Maximum waiting time for inference calls
    MaxWorkers: 8 # Maximum number of concurrent inferences (recommended to be the same as the number of CPU cores)
```



##### 4. Http server 

settings.yaml at 192.168.0.102

```yaml
Server:
    QueueName: "goinfer-synchronous-asynchronous-queue_102"
```



settings.yaml at 192.168.0.103 

```yaml
Server:
    QueueName: "goinfer-synchronous-asynchronous-queue_103"
```

Start command

```bash
build/go-embedding http
```



##### 5. Dispatcher server

###### 192.168.0.104 

settings.yaml

```yaml
Server:
    QueueName: "goinfer-synchronous-asynchronous-queue_102"
```

Start command

```bash
build/go-embedding server 0
```



###### 192.168.0.105 

settings.yaml

```yaml
Server:
    QueueName: "goinfer-synchronous-asynchronous-queue_102"
```

Start command

```bash
build/go-embedding server 1
```



###### 192.168.0.106

settings.yaml

```yaml
Server:
    QueueName: "goinfer-synchronous-asynchronous-queue_103"
```

Start command

```bash
build/go-embedding server 0
```



###### 192.168.0.107

settings.yaml

```yaml
Server:
    QueueName: "goinfer-synchronous-asynchronous-queue_103"
```

Start command

```bash
build/go-embedding server 1
```



##### 6. Testing

```bash
python3 test_api "192.168.0.100" mobile
```
