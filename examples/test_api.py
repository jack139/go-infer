# coding:utf-8

import sys, urllib3, json, base64, time, hashlib
from datetime import datetime
from utils import sm2

urllib3.disable_warnings()

# 生成参数字符串
def gen_param_str(param1):
    param = param1.copy()
    name_list = sorted(param.keys())
    if 'data' in name_list: # data 按 key 排序, 中文不进行性转义，与go保持一致
        param['data'] = json.dumps(param['data'], sort_keys=True, ensure_ascii=False, separators=(',', ':'))
    return '&'.join(['%s=%s'%(str(i), str(param[i])) for i in name_list if str(param[i])!=''])


if __name__ == '__main__':
    if len(sys.argv)<3:
        print("usage: python3 %s <host> <api_path>" % sys.argv[0])
        sys.exit(2)

    hostname = sys.argv[1]
    api_path = sys.argv[2]

    with open("test.jpg", 'rb') as f:
        img_data = f.read()

    body = {
        #'version'  : '1',
        #'signType' : 'SHA256', 
        'signType' : 'SM2',
        #'signType' : 'plain',
        'encType'  : 'plain',
        'data'     : {
            'image'    : "",
            'text'     : "测试测试",
        }
    }

    if api_path in ["mobile", "facedet"]:
        body['data']['image'] = base64.b64encode(img_data).decode('utf-8')

    appid = '3EA25569454745D01219080B779F021F'
    unixtime = int(time.time())
    body['timestamp'] = unixtime
    body['appId'] = appid

    param_str = gen_param_str(body)
    sign_str = '%s&key=%s' % (param_str, '41DF0E6AE27B5282C07EF5124642A352')

    #print(sign_str)

    if body['signType'] == 'SHA256':
        signature_str =  base64.b64encode(hashlib.sha256(sign_str.encode('utf-8')).hexdigest().encode('utf-8')).decode('utf-8')
    elif body['signType'] == 'SM2':
        signature_str = sm2.SM2withSM3_sign_base64(sign_str)
    else: # plain
        signature_str = ""

    body['signData'] = signature_str

    body = json.dumps(body)
    #print(body)

    pool = urllib3.PoolManager(num_pools=2, timeout=180, retries=False)

    host = 'http://%s:5000'%hostname
    

    url = host+'/api/'+api_path

    start_time = datetime.now()
    r = pool.urlopen('POST', url, body=body)
    print('[Time taken: {!s}]'.format(datetime.now() - start_time))

    print(r.status)
    if r.status==200:
        print(json.loads(r.data.decode('utf-8')))
    else:
        print(r.data)
