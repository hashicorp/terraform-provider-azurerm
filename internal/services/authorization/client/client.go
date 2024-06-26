// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck // nolint: staticcheck

	// To swap sdk for `azurerm_role_definition` without changing API version
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentscheduleinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentschedulerequests"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityscheduleinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedulerequests"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicyassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RoleAssignmentsClient                  *authorization.RoleAssignmentsClient
	RoleAssignmentScheduleRequestClient    *roleassignmentschedulerequests.RoleAssignmentScheduleRequestsClient
	RoleAssignmentScheduleInstancesClient  *roleassignmentscheduleinstances.RoleAssignmentScheduleInstancesClient
	RoleAssignmentSchedulesClient          *roleassignmentschedules.RoleAssignmentSchedulesClient
	RoleEligibilityScheduleRequestClient   *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestsClient
	RoleEligibilityScheduleInstancesClient *roleeligibilityscheduleinstances.RoleEligibilityScheduleInstancesClient
	RoleEligibilitySchedulesClient         *roleeligibilityschedules.RoleEligibilitySchedulesClient
	RoleManagementPoliciesClient           *rolemanagementpolicies.RoleManagementPoliciesClient
	RoleManagementPolicyAssignmentsClient  *rolemanagementpolicyassignments.RoleManagementPolicyAssignmentsClient
	ScopedRoleAssignmentsClient            *roleassignments.RoleAssignmentsClient
	ScopedRoleDefinitionsClient            *roledefinitions.RoleDefinitionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	roleAssignmentsClient := authorization.NewRoleAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	roleAssignmentScheduleRequestsClient, err := roleassignmentschedulerequests.NewRoleAssignmentScheduleRequestsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating roleAssignmentScheduleRequestsClient: %+v", err)
	}

	o.Configure(roleAssignmentScheduleRequestsClient.Client, o.Authorizers.ResourceManager)

	roleAssignmentScheduleInstancesClient, err := roleassignmentscheduleinstances.NewRoleAssignmentScheduleInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating roleAssignmentScheduleInstancesClient: %+v", err)
	}
	o.Configure(roleAssignmentScheduleInstancesClient.Client, o.Authorizers.ResourceManager)

	roleAssignmentSchedulesClient, err := roleassignmentschedules.NewRoleAssignmentSchedulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating roleAssignmentSchedulesClient: %+v", err)
	}
	o.Configure(roleAssignmentSchedulesClient.Client, o.Authorizers.ResourceManager)

	roleEligibilityScheduleRequestClient, err := roleeligibilityschedulerequests.NewRoleEligibilityScheduleRequestsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating roleEligibilityScheduleRequestClient: %+v", err)
	}
	o.Configure(roleEligibilityScheduleRequestClient.Client, o.Authorizers.ResourceManager)

	roleEligibilityScheduleInstancesClient, err := roleeligibilityscheduleinstances.NewRoleEligibilityScheduleInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating roleEligibilityScheduleInstancesClient: %+v", err)
	}
	o.Configure(roleEligibilityScheduleInstancesClient.Client, o.Authorizers.ResourceManager)

	roleEligibilitySchedulesClient, err := roleeligibilityschedules.NewRoleEligibilitySchedulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating roleEligibilitySchedulesClient: %+v", err)
	}
	o.Configure(roleEligibilitySchedulesClient.Client, o.Authorizers.ResourceManager)

	roleManagementPoliciesClient, err := rolemanagementpolicies.NewRoleManagementPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating roleManagementPoliciesClient: %+v", err)
	}
	o.Configure(roleManagementPoliciesClient.Client, o.Authorizers.ResourceManager)

	roleManagementPolicyAssignmentClient, err := rolemanagementpolicyassignments.NewRoleManagementPolicyAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating roleManagementPolicyAssignmentClient: %+v", err)
	}
	o.Configure(roleManagementPolicyAssignmentClient.Client, o.Authorizers.ResourceManager)

	scopedRoleAssignmentsClient, err := roleassignments.NewRoleAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Role Assignment Client:  %+v", err)
	}
	o.Configure(scopedRoleAssignmentsClient.Client, o.Authorizers.ResourceManager)

	scopedRoleDefinitionsClient, err := roledefinitions.NewRoleDefinitionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Role Definition Client:  %+v", err)
	}
	o.Configure(scopedRoleDefinitionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		RoleAssignmentsClient:                  &roleAssignmentsClient,
		RoleAssignmentScheduleRequestClient:    roleAssignmentScheduleRequestsClient,
		RoleAssignmentScheduleInstancesClient:  roleAssignmentScheduleInstancesClient,
		RoleAssignmentSchedulesClient:          roleAssignmentSchedulesClient,
		RoleEligibilityScheduleRequestClient:   roleEligibilityScheduleRequestClient,
		RoleEligibilityScheduleInstancesClient: roleEligibilityScheduleInstancesClient,
		RoleEligibilitySchedulesClient:         roleEligibilitySchedulesClient,
		RoleManagementPoliciesClient:           roleManagementPoliciesClient,
		RoleManagementPolicyAssignmentsClient:  roleManagementPolicyAssignmentClient,
		ScopedRoleAssignmentsClient:            scopedRoleAssignmentsClient,
		ScopedRoleDefinitionsClient:            scopedRoleDefinitionsClient,
	}, nil
}
