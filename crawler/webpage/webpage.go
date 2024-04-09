package webpage

type WebPage struct {
	URL           string   `json:"url"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Words         []Word   `json:"words"`
	Links         []string `json:"links"`
	Abstract      string   `json:"abstract"`
	OutboundLinks []string `json:"outbound_links"`
}

type Word struct {
	Word     string `json:"word"`
	Score    int    `json:"score"`
	Position int    `json:"position"`
}
