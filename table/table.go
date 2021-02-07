package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"sort"
	"./queue"

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

type Table struct {
	TeamName string
	Points int
}

var matchesUrl string

func init() {
	matchesUrl = os.Getenv("MATCH_URL")
}

func loadMatches() []byte {
	response, err := http.Get(matchesUrl + "/matches")
	if err != nil {
		fmt.Println("Erro de HTTP")
	}
	data, _ := ioutil.ReadAll(response.Body)

	var matches Matches
	json.Unmarshal(data, &matches)

	fmt.Println(string(data))
	return data
}

func processTable(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	data := loadMatches()

	var matches Matches
	json.Unmarshal(data, &matches)

	var table []Table

	for _, v := range matches.Matches {
		if v.TeamHome.Goals > v.TeamAway.Goals {
			//fmt.Printf(v.TeamHome.Team)
			table = append(table, Table{v.TeamHome.Team, 3})
			table = append(table, Table{v.TeamAway.Team, 0})
		} else if v.TeamHome.Goals < v.TeamAway.Goals {
			table = append(table, Table{v.TeamAway.Team, 3})
			table = append(table, Table{v.TeamHome.Team, 0})
		} else {
			table = append(table, Table{v.TeamAway.Team, 1})
			table = append(table, Table{v.TeamHome.Team, 1})
		}
	}

	sort.SliceStable(table, func(i, j int) bool {
    return table[i].Points > table[j].Points
	})

	fmt.Printf("%v", table)

	t := template.Must(template.ParseFiles("templates/table.html"))
	t.Execute(w, table)

	data, _ = json.Marshal(table)
	fmt.Println(string(data))

	connection := queue.Connect()
	queue.Notify(data, "table_ex", "", connection)
}

func finish(w http.ResponseWriter, r *http.Request) {
	match := loadMatches()
	data, _ := json.Marshal(match)

	fmt.Println(string(data))
	
	connection := queue.Connect()
	queue.Notify(data, "table_ex", "", connection)
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("OK!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", test)
	r.HandleFunc("/finish", finish)
	r.HandleFunc("/test", test)
	r.HandleFunc("/table", processTable)
	http.ListenAndServe(":3000", r)
}