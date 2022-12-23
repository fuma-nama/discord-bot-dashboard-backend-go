package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type User struct {
	Id string `json:"id"`
}

type OAuth2Config struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	Scope        string
}

var endpoint = "https://discord.com/api/v10"
var client = &http.Client{}

func CheckToken(accessToken string) bool {
	req, err := request("HEAD", "/oauth2/@me", accessToken)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)

	return err == nil && res.StatusCode == 200
}

func GetToken(options OAuth2Config, callbackUrl string, code string) (*TokenResponse, error) {
	body := url.Values{
		"client_id":     {options.ClientId},
		"client_secret": {options.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {callbackUrl},
	}

	res, err := http.Post(endpoint+"/oauth2/token", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(body.Encode())))

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("failed to exchange token")
	}

	raw, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data TokenResponse
	err = json.Unmarshal(raw, &data)

	return &data, err
}

func RevokeToken(options OAuth2Config, accessToken string) error {
	body := url.Values{
		"client_id":     {options.ClientId},
		"client_secret": {options.ClientSecret},
		"token":         {accessToken},
	}

	res, err := http.Post("https://discord.com/api/oauth2/token/revoke", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(body.Encode())))

	if err != nil || res.StatusCode != 200 {
		return errors.New("failed to revoke token")
	}

	return nil
}

func GetUser(accessToken string) (user *User, err error) {
	req, err := request("GET", "/users/@me", accessToken)
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var data User
	err = json.Unmarshal(raw, &data)

	return &data, nil
}

func request(method, url string, accessToken string) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint+url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	return req, nil
}
