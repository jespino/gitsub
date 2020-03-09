package delivery

import "fmt"

type ConsoleDelivery struct{}

func (d ConsoleDelivery) Notify(repo *Repo, issue *Issue) {
	fmt.Printf("  New issue detected in the repository %s/%s: %s (%s)\n", repo.Owner, repo.Name, issue.Title, issue.URL)
}
