package common

import (
    "crypto/rand"
    "crypto/md5"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "fmt"
    "io"
    "github.com/robfig/revel"
)

const (
    RANDOM_SOURCE = ":;|,<.>?[{]}-_=+~!@#$%^&*()0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func NewRandom(length int) []byte {
    bytes := make([]byte, length)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = RANDOM_SOURCE[b % byte(len(RANDOM_SOURCE))]
    }
    return bytes
}

func NewRandomString(length int) string {
    return string(NewRandom(length))
}

// generate 32byte md5 string
func MD5Sum(source string) string {
    crypt := md5.New()
    io.WriteString(crypt, source)
    return fmt.Sprintf("%x", crypt.Sum(nil))
}

// iv is saved in cipher string[:aes.BlockSize]
func AesEncrypt(source, key []byte) string {
    block, err := aes.NewCipher(key)
    if nil != err {
        revel.ERROR.Printf("failed to create aes cipher: %s", err)
        panic(err)
    }

    sourceStr := EncodeBase64(source)
    cipherStr := make([]byte, aes.BlockSize + len(sourceStr))
    _, err = rand.Read(cipherStr[:aes.BlockSize])
    if nil != err {
        revel.ERROR.Printf("failed to generate random bytes: %s", err)
        panic(err)
    }
    iv := cipherStr[:aes.BlockSize]
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(cipherStr[aes.BlockSize:], []byte(sourceStr))

    return string(cipherStr)
}

func AesDecrypt(source, key []byte) string {
    block, err := aes.NewCipher(key)
    if nil != err {
        revel.ERROR.Printf("failed to create aes cipher: %s", err)
        panic(err)
    }

    if len(source) < aes.BlockSize {
        revel.ERROR.Printf("decrypt error: cipher string too short %s", source)
        return ""
    }

    iv := source[:aes.BlockSize]
    cipherStr := source[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(cipherStr, cipherStr)

    return string(DecodeBase64(string(cipherStr)))
}

func EncodeBase64(source []byte) string {
    return base64.StdEncoding.EncodeToString(source)
}

func DecodeBase64(source string) []byte {
    data, err := base64.StdEncoding.DecodeString(source)
    if nil != err {
        revel.ERROR.Printf("failed to decode %s with base64", source)
        return []byte("")
    }
    return data
}

