package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Event struct {
	Type string `json:"type"`
	Repo Repo `json:"repo"`
	Payload Payload `json:"payload"`
}

type Repo struct {
	Name string `json:"name"`	
}

type Payload struct{
	Action string `json:"action"`
}

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}
	username := os.Args[1]
	ghApiUrl := "https://api.github.com"
	req,err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/events", ghApiUrl, username), nil)
	if err != nil {
		log.Fatal(err)
	}
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d\n", res.StatusCode)
	}

	var events []Event

	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		log.Fatal(err)
	}
	// print pushed commits summary
	count_commits := make(map[string]int,len(events))
	for i:=0; i< len(events);i++{
		if events[i].Type == "PushEvent"{
			count_commits[events[i].Repo.Name] += 1
		}
	}
	fmt.Println("Output:")
	for k, v := range count_commits{
		fmt.Printf("  Pushed %d commits to %s repo\n",v,k)
	}

	// print pushed pull request  summary
	count_prs := make(map[string]int,len(events))
	for i:=0; i< len(events);i++{
		if events[i].Type == "PullRequestEvent"{
			count_prs[events[i].Repo.Name] += 1
		}
	}
	for k, v := range count_prs{
		fmt.Printf("  Contributed to %d pr %s repo\n",v,k)
	}
	
  // print contributed comments summary
	countComments := make(map[string]int,len(events))
	for i:=0; i< len(events);i++{
		if events[i].Type == "IssueCommentEvent"{
			countComments[events[i].Repo.Name] += 1
		}
	}
	for k, v := range countComments{
		fmt.Printf("  Added %d comments in %s repo\n",v,k)
	}
}

func usage() {
	fmt.Println(os.Args[0], "<username>")
}
