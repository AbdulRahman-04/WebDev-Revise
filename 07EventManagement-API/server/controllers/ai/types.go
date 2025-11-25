package ai

// For Event Description AI
type AiEventInput struct {
	EventName   string `json:"eventname" binding:"required"`
	EventType   string `json:"eventtype" binding:"required"`
	Attendence  int    `json:"attendence" binding:"required"`
	EventDesc   string `json:"eventdesc"`
	Location    string `json:"location" binding:"required"`
}

// For Function Description AI
type AiFunctionInput struct {
	FuncName   string `json:"funcname" binding:"required"`
	FuncType   string `json:"functype" binding:"required"`
	Guests     int    `json:"guests" binding:"required"`
	FuncDesc   string `json:"funcdesc"`
	Location   string `json:"location" binding:"required"`
}
