# -*- coding: utf-8 -*-

# 调用 C 实现的 SM3 摘要

import os, sys
import ctypes
import time
from functools import wraps

__all__ = ['SM3Context', 'lib_sm3', 'sm3_starts', 'sm3_update', 'sm3_finish', 'sm3']

c_ubyte_p = ctypes.POINTER(ctypes.c_ubyte)

class SM3Context(ctypes.Structure):
    _fields_ = [("total", ctypes.c_ulong * 2),
                ("state", ctypes.c_ulong * 8),
                ("buffer", ctypes.c_ubyte * 64),
                ("ipad", ctypes.c_ubyte * 64),
                ("opad", ctypes.c_ubyte * 64)]


lib_sm3 = None
sm3_starts = None
sm3_update = None
sm3_finish = None
sm3 = None

# 只支持 linux 平台
lib_file = "libsm3.so"

try:
    library_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), lib_file)
    if not os.path.exists(library_path):
        library_path = os.path.join(os.getcwd(), lib_file)

    lib_sm3 = ctypes.CDLL(library_path)
except Exception as e:
    lib_sm3 = None


if lib_sm3 is not None:
    sm3_starts = lib_sm3.sm3_starts
    sm3_starts.restype = None
    sm3_starts.argtypes = (ctypes.POINTER(SM3Context), )

    sm3_update = lib_sm3.sm3_update
    sm3_update.restype = None
    sm3_update.argtypes = (ctypes.POINTER(SM3Context), ctypes.POINTER(ctypes.c_ubyte), ctypes.c_int)

    sm3_finish = lib_sm3.sm3_finish
    sm3_finish.restype = None
    sm3_finish.argtypes = (ctypes.POINTER(SM3Context), ctypes.POINTER(ctypes.c_ubyte * 32))

    sm3 = lib_sm3.sm3
    sm3.restype = None
    sm3.argtypes = (ctypes.POINTER(ctypes.c_ubyte), ctypes.c_int, ctypes.POINTER(ctypes.c_ubyte * 32))


def sm3_hash(in_msg):
    sm3_ctx = SM3Context()
    hash_buffer = (ctypes.c_ubyte * 32)()
    sm3_starts(ctypes.POINTER(SM3Context)(sm3_ctx))
    array_type = ctypes.c_ubyte * len(in_msg)
    sm3_update(ctypes.POINTER(SM3Context)(sm3_ctx), array_type(*in_msg), ctypes.c_int(len(in_msg)))
    sm3_finish(ctypes.POINTER(SM3Context)(sm3_ctx), hash_buffer)
    return bytes([i for i in hash_buffer]).hex()


if __name__ == '__main__':
    print(sm3_hash(b'321'))