package nodetype

type EndpointRangeDescription struct {
	EndPort   int64 `json:"endPort"`
	StartPort int64 `json:"startPort"`
}
