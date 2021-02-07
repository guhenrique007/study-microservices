package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
)

type Team struct {
	Team  string `json:"team"`
	Goals string `json:"goals"`
}

type Match struct {
	Uuid     string `json:"uuid"`
	TeamHome Team   `json:"team_home"`
	TeamAway Team   `json:"team_away"`
}

type Matches struct {
	Matches []Match
}

var matchesUrl string

func init() {
	matchesUrl = os.Getenv("MATCH_URL")
}

func loadMatches() []Match {
	response, err := http.Get(matchesUrl + "/matches")
	if err != nil {
		fmt.Println("Erro de HTTP")
	}
	data, _ := ioutil.ReadAll(response.Body)

	var matches Matches
	json.Unmarshal(data, &matches)

	fmt.Println(string(data))
	return matches.Matches
}

func ListMatches(w http.ResponseWriter, r *http.Request) {
	matches := loadMatches()
	t := template.Must(template.ParseFiles("templates/championship.html"))
	t.Execute(w, matches)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", ListMatches)
	http.ListenAndServe(":5000", r)
}

