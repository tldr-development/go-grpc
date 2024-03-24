package inspire

type Inspire struct {
	UUID    string // account uuid
	Status  string // status 발송 여부
	Prompt  string // prompt
	Context string // context
	Message string // response llm 에서 생성한 메시지
	Created string // created at timestamp llm 에서 생성한 시간
	Updated string // updated at timestamp Status가 변경된 시간
	NameKey string
}
