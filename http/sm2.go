package http

import (
	"log"
	"math/big"
	"encoding/base64"
	"crypto/rand"
	"github.com/tjfoc/gmsm/sm2"

	"github.com/jack139/go-infer/helper"
)


var (
	// 用户密钥
	privKey *sm2.PrivateKey
	pubKey *sm2.PublicKey
)

func initSM2(){
	// base64 恢复私钥
	privKey, err := restoreKey(helper.Settings.Api.SM2Private)
	if err!=nil {
		log.Fatal("SM2 private key FAIL")
	}

	// 公钥
	pubKey = &privKey.PublicKey

	log.Printf("D: %x\nX: %x\nY: %x\n", privKey.D, privKey.PublicKey.X, privKey.PublicKey.Y)
}

// 从 base64私钥 恢复密钥对
func restoreKey(privStr string) (*sm2.PrivateKey, error) {
	priv, err  := base64.StdEncoding.DecodeString(privStr)
	if err!=nil {
		return nil, err
	}

	curve := sm2.P256Sm2()
	key := new(sm2.PrivateKey)
	key.PublicKey.Curve = curve
	key.D = new(big.Int).SetBytes(priv)
	key.PublicKey.X, key.PublicKey.Y = curve.ScalarBaseMult(priv)
	return key, nil
}

// SM2签名
func sm2Sign(data []byte) ([]byte, error) {
	// 签名(内部已经做了sm3摘要)
	R, S, err := sm2.Sm2Sign(privKey, data, nil, rand.Reader) 
	if err!=nil {
		return nil, err
	}

	sign := R.Bytes()
	sign = append(sign, S.Bytes()...)

	return sign, nil
}

// SM2签名，返回base64编码
func sm2SignBase64(data []byte) (string, error) {
	sign, err := sm2Sign(data)
	if err!=nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

// SM2验签
func sm2Verify(data []byte, sign []byte) bool {
	if len(sign)<64 {
		return false
	}

	R := new(big.Int).SetBytes(sign[:32]) 
	S := new(big.Int).SetBytes(sign[32:])

	// 验签
	return sm2.Sm2Verify(pubKey, data, nil, R, S)
}

// SM2验签，使用base64编码
func sm2VerifyBase64(data []byte, signBase64 string) bool {
	sign, _  := base64.StdEncoding.DecodeString(signBase64)
	return sm2Verify(data, sign)
}
