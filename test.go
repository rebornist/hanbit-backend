package main

import (
	"crypto"
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
)

func main() {
	ExampleSignPKCS1v15()
}

func Exam1() {
	key := "!@#!@231sadasdasda"
	s := "hello world"

	block, err := aes.NewCipher([]byte(key)) // AES 대칭키 암호화 블록 생성
	if err != nil {
		fmt.Println(err)
		return
	}

	ciphertext := make([]byte, len(s))
	block.Encrypt(ciphertext, []byte(s)) // 평문을 AES 알고리즘으로 암호화
	fmt.Printf("%x\n", ciphertext)

	plaintext := make([]byte, len(s))
	block.Decrypt(plaintext, ciphertext) // AES 알고리즘으로 암호화된 데이터를 평문으로 복호화
	fmt.Println(string(plaintext))
}

func ExampleSignPKCS1v15() {
	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader

	message := []byte("message to be signed")

	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed. This requires
	// that the hash function be collision resistant. SHA-256 is the
	// least-strong hash function that should be used for this at the time
	// of writing (2016).
	hashed := sha256.Sum256(message)

	signature, err := SignPKCS1v15(rng, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}

	fmt.Printf("Signature: %x\n", signature)
}
