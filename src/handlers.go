package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	resty "gopkg.in/resty.v0"
)

type (
	authPayload struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	ldapResponse struct {
		AccessToken      string `json:"access_token,omitempty"`
		RefreshToken     string `json:"refresh_token,omitempty"`
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
	}
)

func checkAuth(w http.ResponseWriter, req *http.Request) {
	var payload authPayload
	var result ldapResponse
	json.NewDecoder(req.Body).Decode(&payload)
	_, err := resty.R().
		SetQueryParams(map[string]string{
			"grant_type": "password",
			"username":   payload.Username,
			"password":   payload.Password,
			"scope":      "read",
			"format":     "json",
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		SetError(&result).
		Post("https://torus-45:jyqgjfawTPj5PrTDPEUI@arel.eisti.fr/oauth/token")

	if err != nil {
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(result)
}
