package domain

type Notification struct {
	LoanID    int32  `json:"ID"`
	Title     string `json:"Title"`
	Status    string `json:"Status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
