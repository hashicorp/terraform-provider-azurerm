package keys

type Trigger struct {
	TimeAfterCreate  *string `json:"timeAfterCreate,omitempty"`
	TimeBeforeExpiry *string `json:"timeBeforeExpiry,omitempty"`
}
