package main

//go:generate go-bindata-assetfs static/... templates/...

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/nbari/violetear"
)

// Model of stuff to render a page
type Model struct {
	IP      string
	Country string
}

// JSONAPIResponse https://developers.google.com/recaptcha/docs/verify
type JSONAPIResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"` // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname    string    `json:"hostname"`     // the hostname of the site where the reCAPTCHA was solved
	ErrorCodes  []int     `json:"error-codes"`  // optional
}

// Templates with functions available to them
var templates = template.New("").Funcs(templateMap)

// Parse all of the bindata templates
func init() {
	for _, path := range AssetNames() {
		bytes, err := Asset(path)
		if err != nil {
			log.Panicf("Unable to parse: path=%s, err=%s", path, err)
		}
		templates.New(path).Parse(string(bytes))
	}
}

// Render a template given a model
func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/index.html", nil)
}

func v2(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/v2.html", nil)
}

func invisible(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/invisible.html", nil)
}

func post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	response := r.FormValue("g-recaptcha-response")
	// did we get a proper recaptcha response? if null, redirect back to sigup page
	if response == "" {
		w.Write([]byte("No recaptha, developer should handle this by printing again the form, etc, go back and test again."))
		return
	}

	secret := "-- secret key --"

	invisible := r.FormValue("invisible")
	if invisible != "" {
		secret = "-- secret key --"
	}

	postURL := "https://www.google.com/recaptcha/api/siteverify"
	remoteip := ""

	postStr := url.Values{"secret": {secret}, "response": {response}, "remoteip": {remoteip}}
	responsePost, err := http.PostForm(postURL, postStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer responsePost.Body.Close()

	body, err := ioutil.ReadAll(responsePost.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var APIResp JSONAPIResponse
	json.Unmarshal(body, &APIResp)
	fmt.Println(APIResp)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(body))
}

func main() {
	router := violetear.New()
	router.LogRequests = true

	// Example route that takes one rest style option
	router.HandleFunc("*", index)
	router.HandleFunc("v2", v2)
	router.HandleFunc("invisible", invisible)
	router.HandleFunc("post", post, "POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
