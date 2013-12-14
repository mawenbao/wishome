package common

import (
    "crypto/rand"
    "crypto/md5"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "encoding/hex"
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
func AesEncrypt(source, key []byte) ([]byte, bool) {
    block, err := aes.NewCipher(key)
    if nil != err {
        revel.ERROR.Printf("failed to create aes cipher: %s", err)
        return nil, false
    }

    sourceStr := EncodeBase64(source)
    cipherStr := make([]byte, aes.BlockSize + len(sourceStr))
    _, err = rand.Read(cipherStr[:aes.BlockSize])
    if nil != err {
        revel.ERROR.Printf("failed to generate random bytes: %s", err)
        return nil, false
    }
    iv := cipherStr[:aes.BlockSize]
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(cipherStr[aes.BlockSize:], []byte(sourceStr))

    return cipherStr, true
}

func AesDecrypt(source, key []byte) ([]byte, bool) {
    block, err := aes.NewCipher(key)
    if nil != err {
        revel.ERROR.Printf("failed to create aes cipher: %s", err)
        return nil, false
    }

    if len(source) < aes.BlockSize {
        revel.ERROR.Printf("decrypt error: cipher string too short %s", source)
        return nil, false
    }

    iv := source[:aes.BlockSize]
    cipherStr := source[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(cipherStr, cipherStr)

    return DecodeBase64(string(cipherStr))
}

func EncodeBase64(source []byte) string {
    return base64.StdEncoding.EncodeToString(source)
}

func DecodeBase64(source string) ([]byte, bool) {
    data, err := base64.StdEncoding.DecodeString(source)
    if nil != err {
        revel.ERROR.Printf("failed to decode %s with base64", source)
        return nil, false
    }
    return data, true
}

func EncodeToHexString(source []byte) string {
    return hex.EncodeToString(source)
}

func DecodeHexString(source string) ([]byte, bool) {
    target, err := hex.DecodeString(source)
    if nil != err {
        revel.ERROR.Printf("failed to decode hex string %s", source)
        return nil, false
    }
    return target, true
}

