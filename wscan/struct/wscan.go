package wscan

type Wscan struct {
	UUID    string // account uuid
	Status  string // status 발송 여부
	Prompt  string // prompt
	Context string // context
	Message string // response llm 에서 생성한 메시지
	Created int64  // created at timestamp llm 에서 생성한 시간
	Updated int64  // updated at timestamp Status가 변경된 시간
	NameKey string
}
