# coding=utf-8

import tensorflow as tf
import os, shutil
import utils

os.environ["CUDA_VISIBLE_DEVICES"] = '0'

# ------- bert_google base
from bert_modeling import BertModel, BertConfig
config_json_path = '../../../../nlp/nlp_model/chinese_bert_L-12_H-768_A-12/bert_config.json'
checkpoint_path = '../../../../nlp/nlp_model/chinese_bert_L-12_H-768_A-12/bert_model.ckpt'

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

utils.init_from_checkpoint(checkpoint_path)

config = tf.ConfigProto()
config.allow_soft_placement = True
config.gpu_options.allow_growth = True

if os.path.exists('outputs/saved-model'):
    shutil.rmtree("outputs/saved-model") 
os.makedirs("outputs/saved-model")

with tf.Session(config=config) as sess:
    sess.run(tf.global_variables_initializer())

    #print('\n'.join([n.name for n in tf.get_default_graph().as_graph_def().node])) # 所有层的名字

    # save_model 输出 , for goland 测试
    builder = tf.saved_model.builder.SavedModelBuilder("outputs/saved-model")
    builder.add_meta_graph_and_variables(sess, [tf.saved_model.tag_constants.TRAINING], clear_devices=True)
    builder.save()  
