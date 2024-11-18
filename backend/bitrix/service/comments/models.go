package comments

type Comment struct {
	ID      string `json:"ID"`
	Comment string `json:"COMMENT"`
}

type CommentsResponse struct {
	Result []Comment `json:"result"`
	Total  int       `json:"total"`
	Time   struct {
		Start            float64 `json:"start"`
		Finish           float64 `json:"finish"`
		Duration         float64 `json:"duration"`
		Processing       float64 `json:"processing"`
		DateStart        string  `json:"date_start"`
		DateFinish       string  `json:"date_finish"`
		OperatingResetAt int     `json:"operating_reset_at"`
		Operating        int     `json:"operating"`
	} `json:"time"`
}
