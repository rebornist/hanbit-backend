package main

import (
	"crypto/aes"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/mixins"
)

type SermonResponse struct {
	ID        uint   `json:"id"`
	UserId    string `json:"user_id"`
	Email     string `json:"author"`
	Title     string `json:"title"`
	PhotoUrl  string `json:"photo_url"`
	Broadcast string `json:"broadcast"`
	PostType  uint   `json:"post_type"`
	// Content   string    `json:"content"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	test := "1616552955034saintmountain8233"
	fmt.Println("0", mixins.Signing(test))
	fmt.Println(mixins.Unsigning(mixins.Signing(test)))

}

// func testSigning(s string) (string, string) {
// 	str1 := mixins.Signing(s)
// 	uDec, _ := base64.StdEncoding.DecodeString(str1)
// 	str2 := mixins.Unsigning(string(uDec))
// 	return str1, str2
// }

func sqlTest(sermons *[]SermonResponse) {

	db := config.ConnectDb()
	// 웹 서비스 정보 중 데이터베이스 정보 추출
	var DB config.Database
	getInfo, err := config.GetServiceInfo("hanbit_database")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(getInfo, &DB)

	tSermon := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["ser"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])
	tImage := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["img"])
	if err := db.
		Table(tSermon).
		Order(fmt.Sprintf("%s.created_at desc", DB.Tables["ser"])).
		Select(fmt.Sprintf(
			"%s.id, %s.user_id, %s.email, %s.title, %s.photo_url, %s.post_type, %s.status, %s.created_at",
			DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["usr"],
			DB.Tables["ser"],
			DB.Tables["img"],
			DB.Tables["ser"],
			// DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["ser"],
		)).
		Joins(fmt.Sprintf("left join %s on %s.user_id = %s.uid left join %s on %s.user_id = %s.user_id", tUser, DB.Tables["ser"], DB.Tables["usr"], tImage, DB.Tables["ser"], DB.Tables["img"])).
		Where(fmt.Sprintf("%s.status = ?", DB.Tables["ser"]), 1).
		Scan(&sermons).Error; err != nil {
		fmt.Println(err)
	}
}

func ImageResize(imagePath string) error {

	// open file
	src, err := imaging.Open(imagePath)
	if err != nil {
		return err
	}

	var divNum float64
	width, height, _ := imageScan(imagePath)
	if width >= height {
		divNum = float64(1980) / float64(width)
	} else {
		divNum = float64(1080) / float64(height)
	}
	width, height = int(float64(width)*divNum), int(float64(height)*divNum)

	src = imaging.Resize(src, width, height, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err = imaging.Save(src, imagePath, imaging.JPEGQuality(90))
	if err != nil {
		return err
	}

	return nil
}

func imageScan(fp string) (width, height int, err error) {
	if reader, err := os.Open(fp); err == nil {
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			return 0, 0, err
		}
		return im.Width, im.Height, nil
	} else {
		return 0, 0, err
	}
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

// func ExampleSignPKCS1v15() {
// 	// crypto/rand.Reader is a good source of entropy for blinding the RSA
// 	// operation.
// 	rng := rand.Reader

// 	message := []byte("message to be signed")

// 	// Only small messages can be signed directly; thus the hash of a
// 	// message, rather than the message itself, is signed. This requires
// 	// that the hash function be collision resistant. SHA-256 is the
// 	// least-strong hash function that should be used for this at the time
// 	// of writing (2016).
// 	hashed := sha256.Sum256(message)

// 	signature, err := SignPKCS1v15(rng, rsaPrivateKey, crypto.SHA256, hashed[:])
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
// 		return
// 	}

// 	fmt.Printf("Signature: %x\n", signature)
// }
