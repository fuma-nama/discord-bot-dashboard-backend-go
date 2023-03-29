package discord

import (
	"encoding/json"
	"io"
	"net/http"
)

type User struct {
	Id string `json:"id"`
}

var endpoint = "https://discord.com/api/v10"
var client = &http.Client{}

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
