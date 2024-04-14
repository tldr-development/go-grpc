package inspire

type Inspire struct {
	UUID    string // account uuid
	Status  string // status 발송 여부
	Prompt  string // prompt
	Context string // context
	Message string // response llm 에서 생성한 메시지
	Created int64  // created at timestamp llm 에서 생성한 시간
	Updated int64  // updated at timestamp Status가 변경된 시간
	NameKey string
}

type Info struct {
	UUID              string // account uuid
	Status            string // 상태 0: active, 1: inactive, 2: deleted, 3: pending, 4: blocked
	NotiPeriod        string // 알림 주기
	MessageLengthType string // 메시지 길이 타입 0: short, 1: middle, 2: long
	MessageType       string // 메시지 타입 0: 상담사, 1: 친구, 2: 편지, 3: 일기, 4: 기타
	Context           string // context
	UserContext       string // user context
	LastMessage       string // 마지막 메시지
	Updated           int64  // updated at timestamp Status가 변경된 시간
	NameKey           string
}
