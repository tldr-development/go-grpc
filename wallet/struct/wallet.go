package wallet

type Wallet struct {
	UUID    string // uuid
	Status  string // status
	Ticket  int64  // ticket
	Created int64  // created at timestamp
	Updated int64  // updated at timestamp
}
