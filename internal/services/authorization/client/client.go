package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentscheduleinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentschedulerequests"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityscheduleinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedulerequests"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RoleAssignmentsClient                  *authorization.RoleAssignmentsClient
	RoleDefinitionsClient                  *authorization.RoleDefinitionsClient
	RoleAssignmentScheduleRequestClient    *roleassignmentschedulerequests.RoleAssignmentScheduleRequestsClient
	RoleAssignmentScheduleInstancesClient  *roleassignmentscheduleinstances.RoleAssignmentScheduleInstancesClient
	RoleEligibilityScheduleRequestClient   *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestsClient
	RoleEligibilityScheduleInstancesClient *roleeligibilityscheduleinstances.RoleEligibilityScheduleInstancesClient
}

func NewClient(o *common.ClientOptions) *Client {
	roleAssignmentsClient := authorization.NewRoleAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	roleDefinitionsClient := authorization.NewRoleDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	roleAssignmentScheduleRequestClient := roleassignmentschedulerequests.NewRoleAssignmentScheduleRequestsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&roleAssignmentScheduleRequestClient.Client, o.ResourceManagerAuthorizer)

	roleAssignmentScheduleInstancesClient := roleassignmentscheduleinstances.NewRoleAssignmentScheduleInstancesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&roleAssignmentScheduleInstancesClient.Client, o.ResourceManagerAuthorizer)

	roleEligibilityScheduleRequestClient := roleeligibilityschedulerequests.NewRoleEligibilityScheduleRequestsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&roleEligibilityScheduleRequestClient.Client, o.ResourceManagerAuthorizer)

	roleEligibilityScheduleInstancesClient := roleeligibilityscheduleinstances.NewRoleEligibilityScheduleInstancesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&roleEligibilityScheduleInstancesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RoleAssignmentsClient:                  &roleAssignmentsClient,
		RoleDefinitionsClient:                  &roleDefinitionsClient,
		RoleAssignmentScheduleRequestClient:    &roleAssignmentScheduleRequestClient,
		RoleAssignmentScheduleInstancesClient:  &roleAssignmentScheduleInstancesClient,
		RoleEligibilityScheduleRequestClient:   &roleEligibilityScheduleRequestClient,
		RoleEligibilityScheduleInstancesClient: &roleEligibilityScheduleInstancesClient,
	}
}
