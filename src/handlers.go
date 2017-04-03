package main

import (
	"encoding/json"
	"net/http"

	"path/filepath"

	"os"

	"github.com/colinmarc/hdfs"
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

	hdfsResponse struct {
		Date string
		Name string
		Size int64
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

	if result.Error != "" {
		redisClient.SET(result.AccessToken, payload.Username)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getFiles(w http.ResponseWriter, req *http.Request) {
	client, err := hdfs.New(*baseHost + ":8020")

	if err != nil {
		glog.Errorln(err.Error())
	}

	params := mux.Vars(req)
	var user string

	user = redisClient.GET(params["id"])

	if user == "" {
		glog.Errorln("User not found")
		http.Error(w, "User not found", 500)
		return
	}
	path := "/user/admin"
	searchDir := filepath.Join(path, "/", user)
	array, e := client.ReadDir(searchDir)

	if e != nil {
		glog.Errorln(e.Error())
		http.Error(w, e.Error(), 500)
		return
	}
	var result []*hdfsResponse
	for _, info := range array {
		date := info.ModTime().UTC().Format("02/01/2006 15:04:05 UTC")
		size := info.Size()
		name := info.Name()

		tmp := &hdfsResponse{
			Date: date,
			Size: size,
			Name: name,
		}
		result = append(result, tmp)
	}

	w.Header().Add("Content-Type", "application/json")
	// w.Write(result)
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
