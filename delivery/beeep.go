package delivery

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

type BeeepDelivery struct{}

func (d BeeepDelivery) Notify(repo *Repo, issue *Issue) {
	beeep.Notify(fmt.Sprintf("New Issue detected on %s/%s", repo.Owner, repo.Name), issue.Title+"\n"+issue.URL, "")
}
