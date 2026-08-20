package model

type Price struct {
	Draw int            `json:"draw"`
	Cost map[string]int `json:"cost"`
}

type Manifest struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Draws  int     `json:"draws"`
	Prices []Price `json:"prices"`
}
