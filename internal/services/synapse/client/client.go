// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	managedvirtualnetwork "github.com/tombuildsstuff/kermit/sdk/synapse/2019-06-01-preview/synapse"
	accesscontrol "github.com/tombuildsstuff/kermit/sdk/synapse/2020-08-01-preview/synapse"
	artifacts "github.com/tombuildsstuff/kermit/sdk/synapse/2021-06-01-preview/synapse"
)

type Client struct {
	FirewallRulesClient                               *synapse.IPFirewallRulesClient
	IntegrationRuntimeAuthKeysClient                  *synapse.IntegrationRuntimeAuthKeysClient
	IntegrationRuntimesClient                         *synapse.IntegrationRuntimesClient
	KeysClient                                        *synapse.KeysClient
	PrivateLinkHubsClient                             *synapse.PrivateLinkHubsClient
	SparkPoolClient                                   *synapse.BigDataPoolsClient
	SqlPoolClient                                     *synapse.SQLPoolsClient
	SqlPoolExtendedBlobAuditingPoliciesClient         *synapse.ExtendedSQLPoolBlobAuditingPoliciesClient
	SqlPoolGeoBackupPoliciesClient                    *synapse.SQLPoolGeoBackupPoliciesClient
	SqlPoolSecurityAlertPolicyClient                  *synapse.SQLPoolSecurityAlertPoliciesClient
	SqlPoolTransparentDataEncryptionClient            *synapse.SQLPoolTransparentDataEncryptionsClient
	SqlPoolVulnerabilityAssessmentsClient             *synapse.SQLPoolVulnerabilityAssessmentsClient
	SQLPoolVulnerabilityAssessmentRuleBaselinesClient *synapse.SQLPoolVulnerabilityAssessmentRuleBaselinesClient
	SQLPoolWorkloadClassifierClient                   *synapse.SQLPoolWorkloadClassifierClient
	SQLPoolWorkloadGroupClient                        *synapse.SQLPoolWorkloadGroupClient
	WorkspaceAadAdminsClient                          *synapse.WorkspaceAadAdminsClient
	WorkspaceClient                                   *synapse.WorkspacesClient
	WorkspaceExtendedBlobAuditingPoliciesClient       *synapse.WorkspaceManagedSQLServerExtendedBlobAuditingPoliciesClient
	WorkspaceManagedIdentitySQLControlSettingsClient  *synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	WorkspaceSecurityAlertPolicyClient                *synapse.WorkspaceManagedSQLServerSecurityAlertPolicyClient
	WorkspaceSQLAadAdminsClient                       *synapse.WorkspaceSQLAadAdminsClient
	WorkspaceVulnerabilityAssessmentsClient           *synapse.WorkspaceManagedSQLServerVulnerabilityAssessmentsClient

	synapseAuthorizer autorest.Authorizer
}

func NewClient(o *common.ClientOptions) *Client {
	firewallRuleClient := synapse.NewIPFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallRuleClient.Client, o.ResourceManagerAuthorizer)

	integrationRuntimeAuthKeysClient := synapse.NewIntegrationRuntimeAuthKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationRuntimeAuthKeysClient.Client, o.ResourceManagerAuthorizer)

	integrationRuntimesClient := synapse.NewIntegrationRuntimesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationRuntimesClient.Client, o.ResourceManagerAuthorizer)

	keysClient := synapse.NewKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&keysClient.Client, o.ResourceManagerAuthorizer)

	privateLinkHubsClient := synapse.NewPrivateLinkHubsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&privateLinkHubsClient.Client, o.ResourceManagerAuthorizer)

	// the service team hopes to rename it to sparkPool, so rename the sdk here
	sparkPoolClient := synapse.NewBigDataPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sparkPoolClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolClient := synapse.NewSQLPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolExtendedBlobAuditingPoliciesClient := synapse.NewExtendedSQLPoolBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolGeoBackupPoliciesClient := synapse.NewSQLPoolGeoBackupPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolGeoBackupPoliciesClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolSecurityAlertPolicyClient := synapse.NewSQLPoolSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolSecurityAlertPolicyClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolTransparentDataEncryptionClient := synapse.NewSQLPoolTransparentDataEncryptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolTransparentDataEncryptionClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolVulnerabilityAssessmentsClient := synapse.NewSQLPoolVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolWorkloadClassifierClient := synapse.NewSQLPoolWorkloadClassifierClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolWorkloadClassifierClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolWorkloadGroupClient := synapse.NewSQLPoolWorkloadGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolWorkloadGroupClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolVulnerabilityAssessmentRuleBaselinesClient := synapse.NewSQLPoolVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolVulnerabilityAssessmentRuleBaselinesClient.Client, o.ResourceManagerAuthorizer)

	workspaceAadAdminsClient := synapse.NewWorkspaceAadAdminsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceAadAdminsClient.Client, o.ResourceManagerAuthorizer)

	workspaceClient := synapse.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceClient.Client, o.ResourceManagerAuthorizer)

	workspaceExtendedBlobAuditingPoliciesClient := synapse.NewWorkspaceManagedSQLServerExtendedBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	workspaceManagedIdentitySQLControlSettingsClient := synapse.NewWorkspaceManagedIdentitySQLControlSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceManagedIdentitySQLControlSettingsClient.Client, o.ResourceManagerAuthorizer)

	workspaceSecurityAlertPolicyClient := synapse.NewWorkspaceManagedSQLServerSecurityAlertPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceSecurityAlertPolicyClient.Client, o.ResourceManagerAuthorizer)

	workspaceSQLAadAdminsClient := synapse.NewWorkspaceSQLAadAdminsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceSQLAadAdminsClient.Client, o.ResourceManagerAuthorizer)

	workspaceVulnerabilityAssessmentsClient := synapse.NewWorkspaceManagedSQLServerVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FirewallRulesClient:                               &firewallRuleClient,
		IntegrationRuntimeAuthKeysClient:                  &integrationRuntimeAuthKeysClient,
		IntegrationRuntimesClient:                         &integrationRuntimesClient,
		KeysClient:                                        &keysClient,
		PrivateLinkHubsClient:                             &privateLinkHubsClient,
		SparkPoolClient:                                   &sparkPoolClient,
		SqlPoolClient:                                     &sqlPoolClient,
		SqlPoolExtendedBlobAuditingPoliciesClient:         &sqlPoolExtendedBlobAuditingPoliciesClient,
		SqlPoolGeoBackupPoliciesClient:                    &sqlPoolGeoBackupPoliciesClient,
		SqlPoolSecurityAlertPolicyClient:                  &sqlPoolSecurityAlertPolicyClient,
		SqlPoolTransparentDataEncryptionClient:            &sqlPoolTransparentDataEncryptionClient,
		SqlPoolVulnerabilityAssessmentsClient:             &sqlPoolVulnerabilityAssessmentsClient,
		SQLPoolVulnerabilityAssessmentRuleBaselinesClient: &sqlPoolVulnerabilityAssessmentRuleBaselinesClient,
		SQLPoolWorkloadClassifierClient:                   &sqlPoolWorkloadClassifierClient,
		SQLPoolWorkloadGroupClient:                        &sqlPoolWorkloadGroupClient,
		WorkspaceAadAdminsClient:                          &workspaceAadAdminsClient,
		WorkspaceClient:                                   &workspaceClient,
		WorkspaceExtendedBlobAuditingPoliciesClient:       &workspaceExtendedBlobAuditingPoliciesClient,
		WorkspaceManagedIdentitySQLControlSettingsClient:  &workspaceManagedIdentitySQLControlSettingsClient,
		WorkspaceSecurityAlertPolicyClient:                &workspaceSecurityAlertPolicyClient,
		WorkspaceSQLAadAdminsClient:                       &workspaceSQLAadAdminsClient,
		WorkspaceVulnerabilityAssessmentsClient:           &workspaceVulnerabilityAssessmentsClient,

		synapseAuthorizer: o.SynapseAuthorizer,
	}
}

func (client Client) RoleDefinitionsClient(workspaceName, synapseEndpointSuffix string) (*accesscontrol.RoleDefinitionsClient, error) {
	if client.synapseAuthorizer == nil {
		return nil, fmt.Errorf("Synapse is not supported in this Azure Environment")
	}
	endpoint := buildEndpoint(workspaceName, synapseEndpointSuffix)
	roleDefinitionsClient := accesscontrol.NewRoleDefinitionsClient(endpoint)
	roleDefinitionsClient.Client.Authorizer = client.synapseAuthorizer
	return &roleDefinitionsClient, nil
}

func (client Client) RoleAssignmentsClient(workspaceName, synapseEndpointSuffix string) (*accesscontrol.RoleAssignmentsClient, error) {
	if client.synapseAuthorizer == nil {
		return nil, fmt.Errorf("Synapse is not supported in this Azure Environment")
	}
	endpoint := buildEndpoint(workspaceName, synapseEndpointSuffix)
	roleAssignmentsClient := accesscontrol.NewRoleAssignmentsClient(endpoint)
	roleAssignmentsClient.Client.Authorizer = client.synapseAuthorizer
	return &roleAssignmentsClient, nil
}

func (client Client) ManagedPrivateEndpointsClient(workspaceName, synapseEndpointSuffix string) (*managedvirtualnetwork.ManagedPrivateEndpointsClient, error) {
	if client.synapseAuthorizer == nil {
		return nil, fmt.Errorf("Synapse is not supported in this Azure Environment")
	}
	endpoint := buildEndpoint(workspaceName, synapseEndpointSuffix)
	managedPrivateEndpointsClient := managedvirtualnetwork.NewManagedPrivateEndpointsClient(endpoint)
	managedPrivateEndpointsClient.Client.Authorizer = client.synapseAuthorizer
	return &managedPrivateEndpointsClient, nil
}

func (client Client) LinkedServiceClient(workspaceName, synapseEndpointSuffix string) (*artifacts.LinkedServiceClient, error) {
	if client.synapseAuthorizer == nil {
		return nil, fmt.Errorf("Synapse is not supported in this Azure Environment")
	}
	endpoint := buildEndpoint(workspaceName, synapseEndpointSuffix)
	linkedServiceClient := artifacts.NewLinkedServiceClient(endpoint)
	linkedServiceClient.Client.Authorizer = client.synapseAuthorizer
	return &linkedServiceClient, nil
}

func buildEndpoint(workspaceName string, synapseEndpointSuffix string) string {
	return fmt.Sprintf("https://%s.%s", workspaceName, synapseEndpointSuffix)
}
