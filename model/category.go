package model

type Category struct {
	Item string `json:"item"`
	Category string `json:"category"`
}

type Categories struct {
	Item string `json:"item"`
	Categories []string `json:"categories"`
}