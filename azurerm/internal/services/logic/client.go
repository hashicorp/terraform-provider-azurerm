package logic

import "github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2016-06-01/logic"

type Client struct {
	WorkflowsClient logic.WorkflowsClient
}
