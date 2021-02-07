package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

func loadData() []byte {
	jsonFile, err := os.Open("matches.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	return data
}

func ListMatches(w http.ResponseWriter, r *http.Request) {
	matches := loadData()
	w.Write([]byte(matches))
}

func ProcessMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := loadData()

	var matches Matches
	json.Unmarshal(data, &matches)

	for _, v := range matches.Matches {
		if v.Uuid == vars["id"] {
			match, _ := json.Marshal(v)
			w.Write([]byte(match))
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/matches", ListMatches)
	r.HandleFunc("/matches/{id}", ProcessMatch)
	http.ListenAndServe(":8081", r)

	fmt.Println(string(loadData()))
}
