// package main

package slackautoinvites

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"appengine"
	"appengine/urlfetch"
)

//Configuration is configuration for slack team
type Configuration struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

func importConfiguration() (string, string) {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration.BaseURL, configuration.Token
}

func setSlackToken(req *http.Request, token string) {
	q := req.URL.Query()
	q.Set("token", token)
	req.URL.RawQuery = q.Encode()
}

func setFormValues(req *http.Request, fname string, lname string, email string) {
	q := req.URL.Query()
	q.Set("first_name", fname)
	q.Set("last_name", lname)
	q.Set("email", email)
	q.Set("set_active", "true")
	q.Set("_attempts", "1")
	req.URL.RawQuery = q.Encode()
}

func sendInvite(r *http.Request, fname string, lname string, email string) (string, error) {
	// client := &http.Client{}

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	baseURL, token := importConfiguration()

	req, _ := http.NewRequest("POST", baseURL, nil)

	setSlackToken(req, token)
	setFormValues(req, fname, lname, email)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
	// TODO - add error checking for body of response
	// success: {"ok":true}
	// failure: {"ok":false,"error":"already_in_team"}
}

func inviteHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("fname") != "" && r.FormValue("lname") != "" && r.FormValue("email") != "" {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		email := r.FormValue("email")

		slackResp, err := sendInvite(r, fname, lname, email)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		fmt.Fprintf(w, slackResp)
	} else {
		fmt.Fprintf(w, "No form values given. Please supply first name, last name, and email.")
	}
}

func init() {
	// func main() {
	http.HandleFunc("/invite", inviteHandler)
}
