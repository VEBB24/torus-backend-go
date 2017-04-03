package main

import (
	"encoding/json"
	"net/http"

	"path/filepath"

	"os"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
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
	glog.Infoln("Process Auth request")
	glog.Infoln(req)

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
		glog.Errorln(err.Error())
		http.Error(w, err.Error(), 500)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getListOfFile(w http.ResponseWriter, req *http.Request) {
	glog.Infoln("Process ListOfFile request")
	glog.Infoln(req)

	params := mux.Vars(req)
	var user string

	user = redisClient.GET(params["id"])

	if user == "" {
		glog.Errorln("User not found")
		http.Error(w, "User not found", 500)
		return
	}

	searchDir := filepath.Join(*basePath, "/", user)
	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		glog.Errorln(err.Error())
		http.Error(w, err.Error(), 500)
	}

	glog.Infoln(fileList)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileList)

}
