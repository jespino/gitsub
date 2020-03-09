package delivery

type Issue struct {
	URL   string
	Title string
}

type Repo struct {
	Owner string
	Name  string
}

type DeliveryService interface {
	Notify(repo *Repo, issue *Issue)
}
