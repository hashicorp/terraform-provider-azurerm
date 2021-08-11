package signalr

type SignalRKeys struct {
	PrimaryConnectionString   *string `json:"primaryConnectionString,omitempty"`
	PrimaryKey                *string `json:"primaryKey,omitempty"`
	SecondaryConnectionString *string `json:"secondaryConnectionString,omitempty"`
	SecondaryKey              *string `json:"secondaryKey,omitempty"`
}
