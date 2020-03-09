package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jespino/gitsub/delivery"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func getLastIssues(client *githubv4.Client, org string, repo string, labels []string, since int64) ([]*delivery.Issue, error) {
	var query struct {
		Repository struct {
			Issues struct {
				Nodes []struct {
					URL   githubv4.String
					Title githubv4.String
				}
			} `graphql:"issues(first: 20, filterBy: {labels: $labels, since: $since}, orderBy: {direction: DESC, field: CREATED_AT})"`
		} `graphql:"repository(name: $repo, owner: $org)"`
	}

	ghLabels := []githubv4.String{}
	for _, label := range labels {
		ghLabels = append(ghLabels, githubv4.String(label))
	}

	variables := map[string]interface{}{
		"org":    githubv4.String(org),
		"repo":   githubv4.String(repo),
		"labels": ghLabels,
		"since":  githubv4.DateTime{Time: time.Unix(since, 0)},
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	issues := []*delivery.Issue{}
	for _, node := range query.Repository.Issues.Nodes {
		issues = append(issues, &delivery.Issue{URL: string(node.URL), Title: string(node.Title)})
	}
	return issues, nil
}

func extractSubscriptions() []string {
	subscriptionsMap := map[string]bool{}
	for _, key := range viper.AllKeys() {
		split := strings.Split(key, ".")
		if len(split) > 0 && split[0] != "config" {
			subscriptionsMap[split[0]] = true
		}
	}

	subscriptions := []string{}
	for subscription := range subscriptionsMap {
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions

}

func main() {
	viper.SetConfigName("gitsub.conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.gitsub")
	viper.AddConfigPath(".")
	viper.SetDefault("config.token", os.Getenv("GITHUB_TOKEN"))
	viper.SetDefault("config.sleep_seconds", 3600)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("config.token")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	deliveryServices := []delivery.DeliveryService{}
	deliveryServices = append(deliveryServices, delivery.ConsoleDelivery{})

	now := time.Now()
	sleepSeconds := viper.GetInt("config.sleep_seconds")
	fmt.Printf("GH Stalker started, checking repositories every %d seconds\n", sleepSeconds)
	subscriptions := extractSubscriptions()
	if len(subscriptions) == 0 {
		fmt.Printf("Not subscriptions found, please review your config file.")
		os.Exit(1)
	}

	for true {
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
		newNow := time.Now()
		for _, sub := range subscriptions {
			owner := viper.GetString(sub + ".owner")
			repo := viper.GetString(sub + ".repo")
			issues, err := getLastIssues(client, owner, repo, viper.GetStringSlice(sub+".labels"), now.Unix())
			if err != nil {
				fmt.Println("Error: ", err.Error())
				os.Exit(1)
			}

			for _, issue := range issues {
				for _, deliveryService := range deliveryServices {
					deliveryService.Notify(&delivery.Repo{Owner: owner, Name: repo}, issue)
				}
			}
		}
		now = newNow
	}
}
