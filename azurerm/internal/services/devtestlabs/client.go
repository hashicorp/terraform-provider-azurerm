package devtestlabs

import "github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"

type Client struct {
	LabsClient            dtl.LabsClient
	PoliciesClient        dtl.PoliciesClient
	VirtualMachinesClient dtl.VirtualMachinesClient
	VirtualNetworksClient dtl.VirtualNetworksClient
}
