package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jespino/gitsub/catalog"
	"github.com/jespino/gitsub/config"
	"github.com/jespino/gitsub/delivery"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch your subscriptions",
	Long:  `Watch your subscriptions and notify you when a new ticket is created`,
	Run:   watchCmdF,
}

func watchCmdF(cmd *cobra.Command, args []string) {
	config, err := config.Read()
	if err != nil {
		fmt.Printf("Fatal error reading config: %s \n", err)
		os.Exit(1)
	}
	catalog, err := catalog.Read()
	if err != nil {
		fmt.Printf("Fatal error reading catalog: %s \n", err)
		os.Exit(1)
	}
	githubToken := config.Main.GithubToken
	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	if len(config.Subscriptions) == 0 && len(config.CustomSubscriptions) == 0 {
		fmt.Printf("Not subscriptions found, please review your config file.")
		os.Exit(1)
	}

	deliveryServices := []delivery.DeliveryService{}
	deliveryServices = append(deliveryServices, delivery.ConsoleDelivery{})

	now := time.Now()
	sleepSeconds := config.Main.SleepSeconds

	fmt.Printf("GitSub watch started, checking repositories every %d seconds\n", sleepSeconds)
	for true {
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
		newNow := time.Now()
		for _, sub := range config.Subscriptions {
			entry, ok := catalog.Catalog[sub.Name]
			if !ok {
				fmt.Printf("Subscription %s not found in your catalog, skiping subscription", sub.Name)
			}
			owner := entry.Org
			repo := entry.Repo
			labels := []string{}
			if sub.HelpWanted {
				labels = append(labels, entry.HelpWanted...)
			}
			if sub.FirstIssue {
				labels = append(labels, entry.FirstIssue...)
			}
			for _, language := range sub.Languages {
				labels = append(labels, entry.Languages[language])
			}
			for _, difficulty := range sub.Difficulties {
				labels = append(labels, entry.Difficulties[difficulty])
			}
			issues, err := getLastIssues(client, owner, repo, labels, now.Unix())
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

		for _, sub := range config.CustomSubscriptions {
			owner := sub.Org
			repo := sub.Repo
			labels := sub.Labels
			issues, err := getLastIssues(client, owner, repo, labels, now.Unix())
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
