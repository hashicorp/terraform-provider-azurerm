package iscsitargets

type Acl struct {
	InitiatorIqn string   `json:"initiatorIqn"`
	MappedLuns   []string `json:"mappedLuns"`
}
