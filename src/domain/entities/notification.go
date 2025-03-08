package domain

type Notification struct {
	LoanID    int32  `json:"loan_id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
