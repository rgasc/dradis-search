package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var config Config

func main() {
	baseUrl := flag.String("baseurl", "", "Dradis base API URL")
	apiKey := flag.String("apikey", "", "Dradis API key")
	term := flag.String("term", "", "Search term")

	flag.Parse()

	if *baseUrl == "" || *apiKey == "" || *term == "" {
		log.Fatal(errors.New("missing required flags"))
	}

	config = Config{
		BaseUrl: *baseUrl + "/%s",
		ApiKey:  *apiKey,
		Term:    *term,
	}

	log.Println("Getting projects from Dradis...")
	projects, err := getProjects()
	if err != nil {
		log.Fatal(err)
	}

	var results []string

	log.Printf("Getting issues and searching for term \"%s\"...", *term)
	for _, project := range projects {
		issues, err := getIssues(project.Id)
		if err != nil {
			log.Fatal(err)
		}

		for _, issue := range issues {
			if strings.Contains(strings.ToLower(issue.Title), strings.ToLower(*term)) {
				results = append(results, fmt.Sprintf(
					config.BaseUrl, fmt.Sprintf("projects/%d/issues/%d (%s)", project.Id, issue.Id, issue.Title),
				))
			}
		}
	}

	if len(results) == 0 {
		log.Println("No results found")
		return
	}

	log.Println("Found the following issues:")
	for i, result := range results {
		fmt.Printf("%d: %s\n", i+1, result)
	}
}

func getResponseBody(r *http.Request) (body []byte, err error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}

	return nil, fmt.Errorf("%d: resource at %s not found", resp.StatusCode, resp.Request.URL)
}

func getProjects() (projects []Project, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(config.BaseUrl, "api/projects"), nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token token=\"%s\"", config.ApiKey))

	res, err := getResponseBody(req)
	if err != nil {
		return
	}

	return projects, json.Unmarshal(res, &projects)
}

func getIssues(projectID int) (issues []Issue, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(config.BaseUrl, "api/issues"), nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token token=\"%s\"", config.ApiKey))
	req.Header.Add("Dradis-Project-Id", strconv.Itoa(projectID))

	res, err := getResponseBody(req)
	if err != nil {
		return
	}

	return issues, json.Unmarshal(res, &issues)
}
