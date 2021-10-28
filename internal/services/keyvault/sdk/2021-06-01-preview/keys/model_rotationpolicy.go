package keys

type RotationPolicy struct {
	Attributes      *KeyRotationPolicyAttributes `json:"attributes,omitempty"`
	LifetimeActions *[]LifetimeAction            `json:"lifetimeActions,omitempty"`
}
