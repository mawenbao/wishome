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
    RANDOM_SOURCE = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    RANDOM_SOURCE_RAW = ":;|,<.>?[{]}-_=+~!@#$%^&*()0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func MakeRandom(length int, randSource string) []byte {
    bytes := make([]byte, length)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = randSource[b % byte(len(randSource))]
    }
    return bytes
}

func NewRawRandom(length int) string {
    return string(MakeRandom(length, RANDOM_SOURCE_RAW))
}

func NewReadableRandom(length int) string {
    return string(MakeRandom(length, RANDOM_SOURCE))
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

