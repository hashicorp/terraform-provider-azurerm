package keys

type LifetimeAction struct {
	Action  *Action  `json:"action,omitempty"`
	Trigger *Trigger `json:"trigger,omitempty"`
}
