# coding=utf-8

# 导出 Google BERT 模型权重，可用于go tf.LoadSavedModel

import tensorflow as tf
import os, shutil

# bert_modeling.py 来自google官方代码 bert/modeling.py
from bert_modeling import BertModel, BertConfig, get_assignment_map_from_checkpoint

os.environ["CUDA_VISIBLE_DEVICES"] = '0'

# ------- google Chinese_BERT-base path
config_json_path = '../../../../nlp/nlp_model/chinese_bert_L-12_H-768_A-12/bert_config.json'
checkpoint_path = '../../../../nlp/nlp_model/chinese_bert_L-12_H-768_A-12/bert_model.ckpt'


# loading weights
def init_from_checkpoint(init_checkpoint, tvars=None, rank=0):
    if not tvars:
        tvars = tf.trainable_variables()
    assignment_map, initialized_variable_names = get_assignment_map_from_checkpoint(tvars, init_checkpoint)
    tf.train.init_from_checkpoint(init_checkpoint, assignment_map)
    if rank == 0:
        # 显示成功加载的权重
        for t in initialized_variable_names:
            if ":0" not in t:
                print("Loading weights success: " + t)


max_seq_length = 512
bert_config = BertConfig.from_json_file(config_json_path)
input_ids = tf.placeholder(tf.int32, shape=[None, max_seq_length], name='input_ids')
segment_ids = tf.placeholder(tf.int32, shape=[None, max_seq_length], name='segment_ids')
input_mask = tf.placeholder(tf.int32, shape=[None, max_seq_length], name='input_mask')
eval_model = BertModel(config=bert_config,
                          is_training=False,
                          input_ids=input_ids,
                          input_mask=input_mask,
                          token_type_ids=segment_ids,
                          use_one_hot_embeddings=False)

init_from_checkpoint(checkpoint_path)

config = tf.ConfigProto()
config.allow_soft_placement = True
config.gpu_options.allow_growth = True

save_model_path = "outputs/saved-model"
if os.path.exists(save_model_path):
    shutil.rmtree(save_model_path) 
os.makedirs(save_model_path)

with tf.Session(config=config) as sess:
    sess.run(tf.global_variables_initializer())

    #print('\n'.join([n.name for n in tf.get_default_graph().as_graph_def().node])) # 所有层的名字

    # save_model 输出 , for goland 测试
    builder = tf.saved_model.builder.SavedModelBuilder(save_model_path)
    builder.add_meta_graph_and_variables(sess, [tf.saved_model.tag_constants.TRAINING], clear_devices=True)
    builder.save()  
