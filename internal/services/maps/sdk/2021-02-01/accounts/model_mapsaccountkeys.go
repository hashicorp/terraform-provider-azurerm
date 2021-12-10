package accounts

type MapsAccountKeys struct {
	PrimaryKey              *string `json:"primaryKey,omitempty"`
	PrimaryKeyLastUpdated   *string `json:"primaryKeyLastUpdated,omitempty"`
	SecondaryKey            *string `json:"secondaryKey,omitempty"`
	SecondaryKeyLastUpdated *string `json:"secondaryKeyLastUpdated,omitempty"`
}
