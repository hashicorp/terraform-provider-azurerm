package mhsmprivateendpointconnections

type ManagedHsmSku struct {
	Family ManagedHsmSkuFamily `json:"family"`
	Name   ManagedHsmSkuName   `json:"name"`
}
