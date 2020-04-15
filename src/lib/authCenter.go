package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"time"
)

const AUTH_URI = "http://ir-auth-inner.duolainc.com"
const AUTH_LOGIN_URI = AUTH_URI + "/api/v1/token/user"
const AUTH_WHOAMI_URI = AUTH_URI + "/api/v1/common/whoami"

var client = &http.Client{
	Timeout: time.Second * 5,
}

type AuthUserTokenResp struct {
	ErrorCode   int    `json:"error_code"`
	ErrorReason string `json:"error_reason"`
	Result      struct {
		AccessToken  interface{} `json:"access_token"`
		ExpiresIn    int         `json:"expires_in"`
		IDToken      string      `json:"id_token"`
		RefreshToken interface{} `json:"refresh_token"`
		TokenType    string      `json:"token_type"`
	} `json:"result"`
}

type UserInfoResp struct {
	ErrorCode   int    `json:"error_code"`
	ErrorReason string `json:"error_reason"`
	Result      struct {
		AccountType string    `json:"account_type"`
		AccountUUID string    `json:"account_uuid"`
		Active      bool      `json:"active"`
		Avatar      string    `json:"avatar"`
		CreatedAt   time.Time `json:"created_at"`
		Departments []struct {
			AccountType string      `json:"account_type"`
			AccountUUID string      `json:"account_uuid"`
			Active      bool        `json:"active"`
			CreatedAt   time.Time   `json:"created_at"`
			Description interface{} `json:"description"`
			ID          int         `json:"id"`
			Name        string      `json:"name"`
			Order       int         `json:"order"`
			Path        string      `json:"path"`
			UpdatedAt   time.Time   `json:"updated_at"`
		} `json:"departments"`
		Email     string        `json:"email"`
		ID        int           `json:"id"`
		IsAdmin   bool          `json:"is_admin"`
		IsLeader  bool          `json:"is_leader"`
		Name      string        `json:"name"`
		NickName  string        `json:"nick_name"`
		Phone     string        `json:"phone"`
		Roles     []interface{} `json:"roles"`
		UpdatedAt time.Time     `json:"updated_at"`
	} `json:"result"`
}

type WechatTokenResp struct {
	AccessToken string `json:"access_token"`
}

func CheckPasswordFromAuth(email string, password string) string {
	requestBody, _ := json.Marshal(map[string]string{
		"username": email,
		"password": password,
	})
	req, err := http.NewRequest("POST", AUTH_LOGIN_URI, bytes.NewBuffer(requestBody))
	if err != nil {
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer res.Body.Close()

	var result AuthUserTokenResp
	if body, err := ioutil.ReadAll(res.Body); err != nil {
		return ""
	} else {
		_ = json.Unmarshal(body, &result)
		return result.Result.IDToken
	}

}

func WhoAmI(idtoken string) *UserInfoResp {
	req, err := http.NewRequest("GET", AUTH_WHOAMI_URI, nil)
	if err != nil {
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+idtoken)
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	var result UserInfoResp
	if body, err := ioutil.ReadAll(res.Body); err != nil {
		return nil
	} else {
		_ = json.Unmarshal(body, &result)
		return &result
	}
}

func AuthGenPassword(password string) string {
	pb := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(pb, bcrypt.DefaultCost)
	return string(hashedPassword)
}

func AuthCheckPassword(password, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	} else {
		return true
	}
}

func GetWechatAppToken() string {
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=ww732d5cac9d0b344f&corpsecret=0-GZQ-T7b2VYbenRZV89debjjKnEwcvXlNKzj2f_D7I"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	var result WechatTokenResp
	if body, err := ioutil.ReadAll(res.Body); err != nil {
		return ""
	} else {
		_ = json.Unmarshal(body, &result)
		return result.AccessToken
	}
}

func SentWechatMarkDown(toUser, msg, token string) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + token
	requestBody, _ := json.Marshal(map[string]interface{}{
		"touser":  toUser,
		"msgtype": "markdown",
		"agentid": 1000032,
		"markdown": map[string]string{
			"content": msg,
		},
		"safe": 0,
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	_, _ = client.Do(req)
}
