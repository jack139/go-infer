# -*- coding: utf-8 -*-

# SM2 验签：

import base64
import binascii
from gmssl import sm2, func
from .libsm3 import sm3

ecc_table = sm2.default_ecc_table

# openssl 1.1.1g
# openssl ecparam -genkey -name SM2 -out priv.key
# openssl ec -in priv.key -pubout -noout -text
#

#pri_key_url_base64 = 'JShsBOJL0RgPAoPttEB1hgtPAvCikOl0V1oTOYL7k5U='

#
# gmssl 需调整： 私钥前加 00, 公钥去掉首字节
private_key = '25286c04e24bd1180f0283edb44075860b4f02f0a290e974575a133982fb9395'
public_key = '4f1df6069a086ac4e1d1c4ad60a3ab26a19ba5fc97a45dedf386c7480dcab18fa745c3a0f6dba6ed6993d0367d9f6b12c06dc01d4079c9eda3f807e21f93edc6'


sm2_crypt = sm2.CryptSM2(public_key=public_key, private_key=private_key)

# SM3WITHSM2 签名规则:  SM2.sign(SM3(Z+MSG)，PrivateKey)
# Z = Hash256(Len(ID) + ID + a + b + xG + yG + xA + yA
def sm3_digest(data_str, pub_key_in_hex_str):
    # sm3withsm2 的 z 值
    z = '0080'+'31323334353637383132333435363738'+ecc_table['a']+ecc_table['b']+ecc_table['g']
    z += pub_key_in_hex_str
    z = binascii.a2b_hex(z)
    Za = sm3.sm3_hash(z)
    M = data_str.encode('utf-8')
    M_ = Za.encode('utf-8') + binascii.b2a_hex(M)
    e = sm3.sm3_hash(binascii.a2b_hex(M_))
    return e

# 加签
def SM2withSM3_sign(data, random_hex_str=None):
    # sm3 摘要
    sm3d = sm3_digest(data, sm2_crypt.public_key)
    sign_data = sm3d.encode('utf-8')
    sign_data = binascii.a2b_hex(sign_data)
    # sm2 加签
    if random_hex_str is None:
        random_hex_str = func.random_hex(sm2_crypt.para_len)
    sign = sm2_crypt.sign(sign_data, random_hex_str) #  16进制
    return sign

# 加签，返回 base64
def SM2withSM3_sign_base64(data, random_hex_str=None):
    sign = SM2withSM3_sign(data, random_hex_str)
    return base64.b64encode(binascii.a2b_hex(sign.encode('utf-8'))).decode('utf-8')

# 验签
def SM2withSM3_verify(sign, data):
    # sm3 摘要
    sm3d = sm3_digest(data, sm2_crypt.public_key)
    sign_data = sm3d.encode('utf-8')
    sign_data = binascii.a2b_hex(sign_data)
    return sm2_crypt.verify(sign, sign_data)

# 验签，签名是 base64
def SM2withSM3_verify_base64(sign_urlbase64, data):
    sign = base64.b64decode(sign_urlbase64)
    sign = binascii.b2a_hex(sign).decode('utf-8')
    return SM2withSM3_verify(sign, data)

