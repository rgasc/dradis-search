package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var config Config

func main() {
	query := flag.String("q", "", "Search query")
	flag.Parse()

	if *query == "" {
		log.Fatal("ERROR: missing search query (-q)")
	}

	if err := loadEnv(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	for _, v := range []string{"BASE_URL", "API_KEY"} {
		val, ok := os.LookupEnv(v)
		if !ok || val == "" {
			log.Fatal("ERROR: Missing required variables in .env file")
		}
	}

	config = Config{
		BaseUrl: os.Getenv("BASE_URL") + "/%s",
		ApiKey:  os.Getenv("API_KEY"),
		Query:   *query,
	}

	log.Println("Getting projects from Dradis...")
	projects, err := getProjects()
	if err != nil {
		log.Fatal(err)
	}

	var results []string

	log.Printf("Getting issues and searching for term \"%s\"...", config.Query)
	for _, project := range projects {
		issues, err := getIssues(project.Id)
		if err != nil {
			log.Fatal(err)
		}

		for _, issue := range issues {
			if strings.Contains(strings.ToLower(issue.Text), strings.ToLower(config.Query)) {
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

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "/dradis-search")
	return filepath.Dir(d)
}

func loadEnv() (err error) {
	envPath := path.Join(rootDir(), ".env")

	if err = godotenv.Load(envPath); err == nil {
		return
	}

	log.Println("Doing a first time setup and storing the variables to .env")

	var baseUrl string
	var apiKey string

	fmt.Println("Specify the Dradis base URL (https://example.com/pro):")
	if _, err = fmt.Scan(&baseUrl); err != nil {
		return
	}

	fmt.Println("Specify the API key (profile -> API Token):")
	if _, err = fmt.Scan(&apiKey); err != nil {
		return
	}

	env, err := godotenv.Unmarshal(fmt.Sprintf("BASE_URL=%s\nAPI_KEY=%s", baseUrl, apiKey))
	if err != nil {
		return
	}

	defer godotenv.Load(envPath)

	return godotenv.Write(env, envPath)
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

	req.Header.Set("Accept", "application/vnd.dradisapi; v=2")
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

	req.Header.Set("Accept", "application/vnd.dradisapi; v=2")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token token=\"%s\"", config.ApiKey))
	req.Header.Add("Dradis-Project-Id", strconv.Itoa(projectID))

	res, err := getResponseBody(req)
	if err != nil {
		return
	}

	return issues, json.Unmarshal(res, &issues)
}
