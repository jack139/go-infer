package mobilenet

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

/*
	below codes taken from
	https://github.com/tensorflow/tensorflow/blob/master/tensorflow/go/example_inception_inference_test.go
*/

// This function constructs a graph of TensorFlow operations which takes as
// input a JPEG-encoded string and returns a tensor suitable as input to the
// inception model.
func constructGraphToNormalizeImage(H, W int32, mean, scale float32, toBGR bool) (graph *tf.Graph, input, output tf.Output, err error) {
	// - The model was trained after with images scaled to scale*scale pixels.
	// - The colors, represented as R, G, B in 1-byte each were converted to
	//   float using (value - Mean)/Scale.

	// - input is a String-Tensor, where the string the JPEG-encoded image. PNG also supported
	// - The inception model takes a 4D tensor of shape
	//   [BatchSize, Height, Width, Colors=3], where each pixel is
	//   represented as a triplet of floats
	// - Apply normalization on each pixel and use ExpandDims to make
	//   this single image be a "batch" of size 1 for ResizeBilinear.

	// toBGR indicated whether changing RGB order to BGR
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	output = op.Div(s,
		op.Sub(s,
			op.ResizeBilinear(s,
				op.ExpandDims(s,
					op.Cast(s,
						op.DecodeJpeg(s, input, op.DecodeJpegChannels(3)), tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), mean)),
		op.Const(s.SubScope("scale"), scale))
	// RGB to BGR
	if toBGR {
		output = op.ReverseV2(s, output, op.Const(s, []int32{-1}))
	}
	graph, err = s.Finalize()
	return graph, input, output, err
}

// Convert the image bytes to a Tensor suitable as input
func makeTensorFromBytes(bytes []byte, H, W int32, mean, scale float32, toBGR bool) (*tf.Tensor, error) {
	// bytes to tensor
	tensor, err := tf.NewTensor(string(bytes))
	if err != nil {
		return nil, err
	}

	// create batch
	graph, input, output, err := constructGraphToNormalizeImage(H, W, mean, scale, toBGR)
	if err != nil {
		return nil, err
	}

	// Execute that graph create the batch of that image
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	batch, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}
	return batch[0], nil
}
