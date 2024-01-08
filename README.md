# Go framework for DL model inference and API deployment

[中文](README_zh.md)

Deep learning models are usually deployed in the cloud and the inference services are provided through APIs. This framework provides the basic architectural components for deploying APIs and achieves several goals:
- The API processing module is decoupled from the model inference module to reduce the risk of network and computing blocking caused by high concurrency.
- The API processing module and model inference module can be deployed as a distributed architecture, and can achieve horizontal expansion.
- The framework implemented using Go language to achieve execution efficiency and simplify deployment and maintenance.
- Custom logic is implemented using callback, hiding common logic. Developers only need to focus on custom logic



Other features:

- Yaml is used for server-side configuration, which can be configured separately during distributed deployment.
- API signature supports SHA256 and SM2 algorithms.
- Model examples:
  - [BERT model inference example](examples/models/embedding)
  - [CNN model inference example](examples/models/mobilenet)
  - [Tensorflow model weight conversion example](examples/export/export_tf_bert.py)
  - [Keras model weight conversion example](examples/export/export_keras_cnn.py)
  - [ONNX model weight inference example](examples/models/facedet)
  - [PyTorch model conversion to ONNX format example](examples/export/pytorch_to_onnx.py)



## Distributed deployment architecture

<img src="doc/arch.png" alt="Distributed deployment architecture" width="500" />



## Development documentation

1. [Development Guide](doc/DEV.md)
2. [API documentation template](doc/API.md)
3. [Local testing](doc/TEST.md)
4. [Tensorflow runtime environment](doc/TF.md)



## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=jack139/go-infer&type=Date)](https://star-history.com/#jack139/go-infer&Date)
