// Package main This app is intented to be go-port of the defunckt's gist library in Ruby
// Currently, uploading single and multiple files are available.
// You can also create secret gists, and both anonymous and user gists.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	version = "0.1.0"
	time    = "2016-08-03-20:30:33"
)

const (
	githubVersion = "v0.3.0"
)

//TODO: A list of clipboard commands with copy and paste support.
//This is intended for adding the gist URLs directly to the user clipboard,
//so that manual copying is not needed.
const (
	xclip   = "xclip -o"
	xsel    = "xsel -o"
	pbcopy  = "pbpaste"
	putclip = "getclip"
)

// Defines different constants used
// GIT_IO_URL is the Github's URL shortner
// API v3 is the current version of GitHub API
const (
	GithubAPIURL = "https://api.github.com/"
	BasePath     = "/api/v3"
)

//User agent defines a custom agent (required by GitHub)
//`token` stores the GITHUB_TOKEN from the env variables
var (
	UserAgent = "gist/#" + githubVersion //Github requires this, else rejects API request
	token     = os.Getenv("GITHUB_TOKEN")
)

// Variables used in `Gist` struct
var (
	publicFlag  bool
	description string
	anonymous   bool
	responseObj map[string]interface{}
)

// GistFile The top-level struct for a gist file
type GistFile struct {
	Content string `json:"content"`
}

// Gist The required structure for POST data for API purposes
type Gist struct {
	Description string              `json:"description"`
	PublicFile  bool                `json:"public"`
	GistFile    map[string]GistFile `json:"files"`
}

//This function loads the GITHUB_TOKEN from a '$HOME/.gist' file
//from the user's home directory.
func loadTokenFromFile() (token string) {
	//get the tokenfile
	file := filepath.Join(os.Getenv("HOME"), ".gist")
	githubToken, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(githubToken))
}

// Defines basic usage when program is run with the help flag
func usage() {
	fmt.Fprintf(os.Stderr, "usage: gister [-p] [-d] [-a] example\n")
	flag.PrintDefaults()
	os.Exit(2)
}

// The main function parses the CLI args. It also checks the files, and
// loads them into an array.
// Then the files are separated into GistFile structs and collectively
// the files are saved in `files` field in the Gist struct.
// A request is then made to the GitHub api - it depends on whether it is
// anonymous gist or not.
// The response recieved is parsed and the Gist URL is printed to STDOUT.
func main() {
	fmt.Println("Build version: ", version, "build time: ", time)
	flag.BoolVar(&publicFlag, "p", false, "Set to false for private gist.")
	flag.BoolVar(&anonymous, "a", false, "Set false if you want the gist for a user")
	flag.StringVar(&description, "d", "This is a gist", "Description for gist.")
	flag.Usage = usage
	flag.Parse()

	fileList := flag.Args()
	if len(fileList) == 0 {
		log.Fatal("Error: No files specified.")
	}

	files := map[string]GistFile{}

	for _, filename := range fileList {
		fmt.Println("Checking file:", filename)
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal("File Error: ", err)
		}
		files[filename] = GistFile{string(content)}
	}

	if description == "" {
		description = strings.Join(fileList, ", ")
	}

	//create a gist from the files array
	gist := Gist{
		description,
		publicFlag,
		files,
	}

	pfile, err := json.Marshal(gist)
	if err != nil {
		log.Fatal("Cannot marshal json: ", err)
	}

	//Check if JSON marshalling succeeds
	fmt.Println("OK")

	b := bytes.NewBuffer(pfile)
	fmt.Println("Uploading...")

	// Send request to API

	req, err := http.NewRequest("POST", "https://api.github.com/gists", b)

	if !anonymous {
		if token == "" {
			token = loadTokenFromFile()
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.SetBasicAuth(token, "x-oauth-basic")
	}

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal("HTTP error: ", err)
	}
	err = json.NewDecoder(response.Body).Decode(&responseObj)
	if err != nil {
		log.Fatal("Response JSON error: ", err)
	}

	fmt.Println("===Gist URL===")
	fmt.Println(responseObj["html_url"])
}
