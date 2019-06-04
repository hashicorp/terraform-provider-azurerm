package policy

import "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"

type Client struct {
	AssignmentsClient    policy.AssignmentsClient
	DefinitionsClient    policy.DefinitionsClient
	SetDefinitionsClient policy.SetDefinitionsClient
}
