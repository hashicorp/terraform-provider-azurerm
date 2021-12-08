package databases

type AccessKeys struct {
	PrimaryKey   *string `json:"primaryKey,omitempty"`
	SecondaryKey *string `json:"secondaryKey,omitempty"`
}
