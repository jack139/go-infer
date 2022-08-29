import numpy as np
import torch
import torchvision.models as models
import onnxruntime as onnxrt


# Load pretrained weights of Resnet50
model = models.resnet50(pretrained=True)

model.eval()

# random input tensor
dummy_input = torch.randn(1, 3, 224, 224)

input_names = [ "actual_input" ]
output_names = [ "output" ]

# convert to onnx model
torch.onnx.export(model, 
                  dummy_input,
                  "outputs/resnet50.onnx",
                  verbose=False,
                  input_names=input_names,
                  output_names=output_names,
                  export_params=True,
                  )


def to_numpy(tensor):
    return tensor.detach().cpu().numpy() if tensor.requires_grad else tensor.cpu().numpy()


# test run onnx model
onnx_session= onnxrt.InferenceSession("outputs/resnet50.onnx")
onnx_inputs= {onnx_session.get_inputs()[0].name: to_numpy(dummy_input)}
onnx_output = onnx_session.run(None, onnx_inputs)
img_label_idx = np.argmax(onnx_output[0])

print(img_label_idx)