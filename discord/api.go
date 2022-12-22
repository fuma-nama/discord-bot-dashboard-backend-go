package discord

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type OAuth2Config struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

var endpoint = "https://discord.com/api/v10"
var client = &http.Client{}

func CheckToken(accessToken string) bool {
	req, err := http.NewRequest("HEAD", endpoint+"/oauth2/@me", nil)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := client.Do(req)

	return err == nil && res.StatusCode == 200
}

func GetToken(options OAuth2Config, code string) (*TokenResponse, error) {
	body, _ := json.Marshal(map[string]string{
		"client_id":     options.ClientId,
		"client_secret": options.ClientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  options.RedirectUrl,
	})

	res, err := http.Post("https://discord.com/api/oauth2/token", "application/x-www-form-urlencoded", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data TokenResponse
	err = json.Unmarshal(raw, &data)

	return &data, err
}