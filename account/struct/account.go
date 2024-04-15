package account

type Account struct {
	UUID    string // uuid
	Status  string // status
	Created string // created at timestamp
	Updated string // updated at timestamp
}

type Platform struct {
	AccountID string // Account uuid id (datastore id)
	Token     string // platform token
	Platform  string // platform name (ex. github, google, kakao)
}
