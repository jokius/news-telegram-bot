package entity

type VkResponse struct {
	VkResult `json:"response"`
}

type VkResult struct {
	Messages []VkMessage `json:"items"`
}

type VkMessage struct {
	ID      uint64 `json:"id"`
	OwnerID int64  `json:"owner_id"`
	Date    int64  `json:"date"`
}
