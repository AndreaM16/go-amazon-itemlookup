package model

type Item struct {
	Item         string `json:"item,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
}

type Items struct {
	Items []Item `json:"items"`
}