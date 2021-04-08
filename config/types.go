package config

// DB 접속 정보
type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Tables   map[string]string
}

// Firebase 접속 정보
type Firebase struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	PrivateKeyAesId         string `json:"private_key_aes_id"`
}

// OAuth Token 정보
type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpriesIn    string `json:"expries_in"`
	TokenType    string `json:"token_type"`
	ProfileUri   string `json:"profile_uri"`
}

// Oauth Profile 정보
type OAuthProfile struct {
	Id           string `json:"id"`
	Nickname     string `json:"nickname"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	MobileE164   string `json:"mobile_e164"`
	Mobile       string `json:"mobile"`
	Birthday     string `json:"birthday"`
	ProfileImage string `json:"profile_image"`
}

// Naver OAuth 2.0 정보
type NaverOAuth struct {
	NaverLoginClientId     string `json:"naver_login_client_id"`
	NaverLoginClientSecret string `json:"naver_login_client_secret"`
	NaverGetTokenUri       string `json:"naver_get_token_uri"`
	NaverGetUserUri        string `json:"naver_get_user_uri"`
}

// Kakao OAuth 정보
type KakaoOAuth struct {
	KakaoLoginClientId     string `json:"kakao_login_client_id"`
	KakaoLoginClientSecret string `json:"kakao_login_client_secret"`
	KakaoGetTokenUri       string `json:"kakao_get_token_uri"`
	KakaoGetUserUri        string `json:"kakao_get_user_uri"`
}

// LetsEncrypt 접속 정보
type Encrypt struct {
	Dir  string `json:"dir"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

// Media 접속 정보
type Media struct {
	Host     string `json:"host"`
	Root     string `json:"root"`
	Ckeditor string `json:"ckeditor"`
	TestRoot string `json:"test_root"`
	TestUser string `json:"test_user"`
}

// service 접속 정보
type Service struct {
	Test string `json:"test"`
}

// 공통 Responst type
type CommonResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}
