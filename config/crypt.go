package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func HanbitEncrypt(b cipher.Block, plaintext []byte) []byte {
	if mod := len(plaintext) % aes.BlockSize; mod != 0 { // 블록 크기의 배수가 되어야함
		padding := make([]byte, aes.BlockSize-mod) // 블록 크기에서 모자라는 부분을
		plaintext = append(plaintext, padding...)  // 채워줌
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext)) // 초기화 벡터 공간(aes.BlockSize)만큼 더 생성
	iv := ciphertext[:aes.BlockSize]                         // 부분 슬라이스로 초기화 벡터 공간을 가져옴
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {  // 랜덤 값을 초기화 벡터에 넣어줌
		fmt.Println(err)
		return nil
	}

	mode := cipher.NewCBCEncrypter(b, iv)                   // 암호화 블록과 초기화 벡터를 넣어서 암호화 블록 모드 인스턴스 생성
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext) // 암호화 블록 모드 인스턴스로
	// 암호화

	return ciphertext
}

func HanbitDecrypt(b cipher.Block, ciphertext []byte) []byte {
	if len(ciphertext)%aes.BlockSize != 0 { // 블록 크기의 배수가 아니면 리턴
		fmt.Println("암호화된 데이터의 길이는 블록 크기의 배수가 되어야합니다.")
		return nil
	}

	iv := ciphertext[:aes.BlockSize]        // 부분 슬라이스로 초기화 벡터 공간을 가져옴
	ciphertext = ciphertext[aes.BlockSize:] // 부분 슬라이스로 암호화된 데이터를 가져옴

	plaintext := make([]byte, len(ciphertext)) // 평문 데이터를 저장할 공간 생성
	mode := cipher.NewCBCDecrypter(b, iv)      // 암호화 블록과 초기화 벡터를 넣어서
	// 복호화 블록 모드 인스턴스 생성
	mode.CryptBlocks(plaintext, ciphertext) // 복호화 블록 모드 인스턴스로 복호화

	return plaintext
}

func EncryptionAESKey(s string) ([]byte, error) {

	b, err := makeCipherBlock()
	if err != nil {
		return nil, err
	}

	plaintext := []byte(s)
	if mod := len(plaintext) % aes.BlockSize; mod != 0 { // 블록 크기의 배수가 되어야함
		padding := make([]byte, aes.BlockSize-mod) // 블록 크기에서 모자라는 부분을
		plaintext = append(plaintext, padding...)  // 채워줌
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext)) // 초기화 벡터 공간(aes.BlockSize)만큼 더 생성
	iv := ciphertext[:aes.BlockSize]                         // 부분 슬라이스로 초기화 벡터 공간을 가져옴
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {  // 랜덤 값을 초기화 벡터에 넣어줌
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(b, iv)                   // 암호화 블록과 초기화 벡터를 넣어서 암호화 블록 모드 인스턴스 생성
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext) // 암호화 블록 모드 인스턴스로 암호화

	return ciphertext, nil
}

func DecryptionAESKey(s string) ([]byte, error) {

	block, err := makeCipherBlock()
	if err != nil {
		return nil, err
	}

	ciphertext := []byte(s)
	if len(ciphertext)%aes.BlockSize != 0 { // 블록 크기의 배수가 아니면 리턴
		return nil, errors.New("암호화된 데이터의 길이는 블록 크기의 배수가 되어야합니다.")
	}

	iv := ciphertext[:aes.BlockSize]        // 부분 슬라이스로 초기화 벡터 공간을 가져옴
	ciphertext = ciphertext[aes.BlockSize:] // 부분 슬라이스로 암호화된 데이터를 가져옴

	plaintext := make([]byte, len(ciphertext)) // 평문 데이터를 저장할 공간 생성
	mode := cipher.NewCBCDecrypter(block, iv)  // 암호화 블록과 초기화 벡터를 넣어서
	// 복호화 블록 모드 인스턴스 생성
	mode.CryptBlocks(plaintext, ciphertext) // 복호화 블록 모드 인스턴스로 복호화

	return plaintext, nil
}

func makeCipherBlock() (cipher.Block, error) {
	// 웹 서비스 정보 중 파이어베이스 정보 추출
	var firebaseInfo Firebase
	getInfo, err := GetServiceInfo("hanbit")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(getInfo, &firebaseInfo)

	b, err := aes.NewCipher([]byte(firebaseInfo.PrivateKeyAesId)) // AES 대칭키 암호화 블록 생성
	if err != nil {
		return nil, err
	}

	return b, nil
}
