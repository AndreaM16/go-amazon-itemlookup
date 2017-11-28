package model

type Item struct {
	Item         string `json:"item,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Title        string `json:"title,omitempty"`
	URL       	 string `json:"url,omitempty"`
	Image        string `json:"image,omitempty"`
}

type Items struct {
	Items []Item `json:"items"`
}