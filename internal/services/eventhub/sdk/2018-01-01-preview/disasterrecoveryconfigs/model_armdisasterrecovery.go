package disasterrecoveryconfigs

type ArmDisasterRecovery struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *ArmDisasterRecoveryProperties `json:"properties,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
