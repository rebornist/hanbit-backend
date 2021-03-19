package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/errorutils"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/mixins"
)

var app = config.InitFirebase()
var ctx = context.Background()
var err error
var req *http.Request
var res *http.Response
var data []byte

// Firebase Auth 초기화
func FirebaseAuthInIt() *auth.Client {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	return client
}

var client = FirebaseAuthInIt()

func FirebaseDeployCookie(code, state, name string) (*http.Cookie, error) {
	var profile config.OAuthProfile
	var token config.OAuthToken
	var cookie *http.Cookie
	var uid string
	var cToken string

	token, err = GetAuthToken(code, state, name)
	if err != nil {
		return nil, err
	}

	profile, err = GetProfile(token.ProfileUri, token.TokenType, token.AccessToken, name)
	if err != nil {
		return nil, err
	}

	uid, err = FirebaseGetUser(profile.Email)
	if err != nil {
		if errorutils.IsNotFound(err) {
			uid, err = FirebaseCreateUser(profile.Email, profile.MobileE164, profile.ProfileImage)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	cToken, err = FirebaseCreateCustomToken(uid)
	if err != nil {
		return nil, err
	}

	err = CheckUser(profile.Email, uid)
	if err != nil {
		return nil, err
	}

	cookie = mixins.CreateCookie("ASESS", cToken, "/api/user")
	return cookie, nil
}

// OAuth 콜백 데이터 전처리
func GetAuthToken(code, state, name string) (config.OAuthToken, error) {
	var address string
	var fToken config.OAuthToken
	var profileUri string
	var reqBody *bytes.Buffer
	var naver config.NaverOAuth
	var kakao config.KakaoOAuth
	contentType := "application/x-www-form-urlencoded;charset=utf-8"

	// 콜백 정보 받아 해당 타입의 주소 생성
	switch name {
	case "naver":
		data, err = getOAuthInfo(name)
		if err != nil {
			return fToken, err
		}
		err := json.Unmarshal(data, &naver)
		if err != nil {
			return fToken, err
		}
		address = fmt.Sprintf(
			"%s?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&state=%s",
			naver.NaverGetTokenUri,
			naver.NaverLoginClientId,
			naver.NaverLoginClientSecret,
			code,
			state,
		)
		profileUri = naver.NaverGetUserUri
	case "kakao":
		data, err = getOAuthInfo(name)
		if err != nil {
			return fToken, err
		}
		err := json.Unmarshal(data, &kakao)
		if err != nil {
			return fToken, err
		}
		address = fmt.Sprintf(
			"%s?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&state=%s",
			kakao.KakaoGetTokenUri,
			kakao.KakaoLoginClientId,
			kakao.KakaoLoginClientSecret,
			code,
			state,
		)
		profileUri = kakao.KakaoGetUserUri
	}

	// 토큰 발급 요청
	reqBody = bytes.NewBufferString(fmt.Sprintf("Post Login %s OAuth", name))
	res, err := http.Post(address, contentType, reqBody)
	if err != nil {
		return fToken, err
	}
	defer res.Body.Close()

	// 결과 출력
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return fToken, err
	}
	json.Unmarshal(data, &fToken)
	fToken.ProfileUri = profileUri

	return fToken, nil
}

// OAuth 프로필 받기
func GetProfile(url, authType, token, name string) (config.OAuthProfile, error) {

	var profile config.OAuthProfile
	var respData map[string]interface{}

	// Request 프로필 발급 요청
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return profile, err
	}

	value := fmt.Sprintf("%s %s", authType, token)

	// Header 추가
	req.Header.Add("Authorization", value)
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=utf-8")

	// Client객체에서 Request 실행
	httpClient := &http.Client{}
	res, err = httpClient.Do(req)
	if err != nil {
		return profile, err
	}
	defer res.Body.Close()

	// 결과 출력
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return profile, err
	}
	json.Unmarshal(data, &respData)

	// 결과 값 중 Profile 값 추출 후 전처리
	switch name {
	case "naver":
		data, err = json.Marshal(respData["response"])
		if err != nil {
			return profile, err
		}
	case "kakao":
		data, err = json.Marshal(respData["kakao_account"])
		if err != nil {
			return profile, err
		}
	}

	// 결과 출력
	json.Unmarshal(data, &profile)

	return profile, nil
}

// User 정보 검색
func FirebaseGetUser(email string) (string, error) {
	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	return user.UID, nil
}

// 유저 생성
func FirebaseCreateUser(email, phone, photo string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		PhoneNumber(phone).
		PhotoURL(photo)
	user, err := client.CreateUser(ctx, params)
	if err != nil {
		return "", err
	}
	return user.UID, nil
}

// 커스텀 토큰 발급
func FirebaseCreateCustomToken(uid string) (string, error) {
	var token string
	token, err = client.CustomToken(ctx, uid)
	if err != nil {
		return token, err
	}
	return token, err
}

// Id 토큰 체크
func FirebaseCheckIdToken(idToken string) error {
	_, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return err
	}
	return nil
}

func getOAuthInfo(name string) ([]byte, error) {
	// 웹 서비스 정보 중 파이어베이스 정보 추출
	getInfo, err := config.GetServiceInfo(name)
	if err != nil {
		return []byte{}, err
	}
	return getInfo, nil
}
