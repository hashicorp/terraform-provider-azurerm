package eventhubs

type EventhubProperties struct {
	CaptureDescription     *CaptureDescription `json:"captureDescription,omitempty"`
	CreatedAt              *string             `json:"createdAt,omitempty"`
	MessageRetentionInDays *int64              `json:"messageRetentionInDays,omitempty"`
	PartitionCount         *int64              `json:"partitionCount,omitempty"`
	PartitionIds           *[]string           `json:"partitionIds,omitempty"`
	Status                 *EntityStatus       `json:"status,omitempty"`
	UpdatedAt              *string             `json:"updatedAt,omitempty"`
}
