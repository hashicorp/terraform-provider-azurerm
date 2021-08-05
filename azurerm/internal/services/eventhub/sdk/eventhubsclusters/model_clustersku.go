package eventhubsclusters

type ClusterSku struct {
	Capacity *int64         `json:"capacity,omitempty"`
	Name     ClusterSkuName `json:"name"`
}
