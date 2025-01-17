package domain

type ItemTypeDto struct {
	Name  string  `json:"name"`
	Type  int     `json:"type"`
	Qty   int     `json:"qty"`
	Image string  `json:"image"`
}

type ItemDto struct {
	Type  int     `json:"type"`
	Qty   int     `json:"qty"`
}
