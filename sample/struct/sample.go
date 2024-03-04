package sample

type Test struct {
	Key   string
	Value string
}

type Account struct {
	UUID          string // uuid
	DeviceAppUUID string // device app uuid
	CreatedAt     int64  // created at
}

type Platform struct {
	UUID     string // Account uuid
	Token    string // platform token
	Platform string // platform name (ex. github, google, kakao)
}

type Profile struct {
	UUID   string // Account uuid
	Name   string // name
	Desc   string // description
	ImgURL string // img url
}
