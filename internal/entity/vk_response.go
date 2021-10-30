package entity

type VkResult struct {
	Messages []VkMessage `json:"items"`
}

type VkMessage struct {
	Date        int            `json:"date"`
	Text        string         `json:"text"`
	Attachments []VkAttachment `json:"attachments"`
}

type VkAttachment struct {
	Type  string  `json:"type"`
	Photo VkPhoto `json:"photo"`
}

type VkPhoto struct {
	Sizes []VkPhotoSize `json:"sizes"`
}

type VkPhotoSize struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	URL    int `json:"url"`
}
