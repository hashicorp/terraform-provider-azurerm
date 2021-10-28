package keys

type KeyAttributes struct {
	Created       *int64                 `json:"created,omitempty"`
	Enabled       *bool                  `json:"enabled,omitempty"`
	Exp           *int64                 `json:"exp,omitempty"`
	Exportable    *bool                  `json:"exportable,omitempty"`
	Nbf           *int64                 `json:"nbf,omitempty"`
	RecoveryLevel *DeletionRecoveryLevel `json:"recoveryLevel,omitempty"`
	Updated       *int64                 `json:"updated,omitempty"`
}
