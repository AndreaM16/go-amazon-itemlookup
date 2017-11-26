package model

type Amazon struct {
	Item Item `json:"item"`
	Manufacturer Manufacturer `json:"manufacturer"`
	Reviews Reviews `json:"review"`
	Categories Categories `json:"categories"`
}