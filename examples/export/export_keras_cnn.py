# coding=utf-8

# 导出 Keras mobilenet 模型权重，可用于go tf.LoadSavedModel


import os, shutil
from keras import backend as K
import tensorflow as tf
from keras.applications import MobileNetV2

os.environ["CUDA_VISIBLE_DEVICES"] = '0'

# ------- Keras applications mobilenet_v2
mobilenet = MobileNetV2(include_top=False, input_shape=(224,224,3), weights='imagenet')
#mobilenet.summary()


save_model_path = "outputs/saved-model_keras"
if os.path.exists(save_model_path):
    shutil.rmtree(save_model_path) 
os.makedirs(save_model_path)


#print('\n'.join([n.name for n in tf.get_default_graph().as_graph_def().node])) # 所有层的名字

# save_model 输出 , for goland 测试
builder = tf.saved_model.builder.SavedModelBuilder(save_model_path)
builder.add_meta_graph_and_variables(K.get_session(), [tf.saved_model.tag_constants.TRAINING], clear_devices=True)
builder.save()  
