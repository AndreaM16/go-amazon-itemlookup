package model

type Review struct {
	Item string `json:"item"`
	Stars float64 `json:"stars"`
	Content string `json:"content"`
	Sentiment uint8 `json:"sentiment"`
	Date string `json:"date"`
}

type Reviews struct {
	Item string `json:"item"`
	Reviews []Review `json:"reviews"`
}