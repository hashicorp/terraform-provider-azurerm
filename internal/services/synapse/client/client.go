// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2020-08-01-preview/roleassignments"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2020-08-01-preview/synapseroledefinitions"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/linkedservices"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/ipfirewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/keys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/privatelinkhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	// Resource Manager
	WorkspacesClient *workspaces.WorkspacesClient

	// Data Plane
	ManagedPrivateEndpointsClient *managedprivateendpoints.ManagedPrivateEndpointsClient

	// TODO: Migrate to go-azure-sdk
	FirewallRulesClient                               *ipfirewallrules.IPFirewallRulesClient
	IntegrationRuntimeAuthKeysClient                  *synapse.IntegrationRuntimeAuthKeysClient
	IntegrationRuntimesClient                         *synapse.IntegrationRuntimesClient
	KeysClient                                        *keys.KeysClient
	PrivateLinkHubsClient                             *privatelinkhubs.PrivateLinkHubsClient
	SparkPoolClient                                   *bigdatapools.BigDataPoolsClient
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
	WorkspaceAzureADOnlyAuthenticationsClient         *synapse.AzureADOnlyAuthenticationsClient
	WorkspaceClient                                   *synapse.WorkspacesClient
	WorkspaceExtendedBlobAuditingPoliciesClient       *synapse.WorkspaceManagedSQLServerExtendedBlobAuditingPoliciesClient
	WorkspaceManagedIdentitySQLControlSettingsClient  *synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	WorkspaceSecurityAlertPolicyClient                *synapse.WorkspaceManagedSQLServerSecurityAlertPolicyClient
	WorkspaceSQLAadAdminsClient                       *synapse.WorkspaceSQLAadAdminsClient
	WorkspaceVulnerabilityAssessmentsClient           *synapse.WorkspaceManagedSQLServerVulnerabilityAssessmentsClient

	// Data Plane (go-azure-sdk) — configured base clients, cloned per workspace endpoint
	roleAssignmentsClient        *roleassignments.RoleAssignmentsClient
	synapseRoleDefinitionsClient *synapseroledefinitions.SynapseRoleDefinitionsClient
	linkedServiceClient          *linkedservices.LinkedServicesClient
	synapseEnabled               bool
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// Resource Manager
	workspacesClient, err := workspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces Client: %+v", err)
	}
	o.Configure(workspacesClient.Client, o.Authorizers.ResourceManager)

	// Data Plane
	managedPrivateEndpointsClient, err := managedprivateendpoints.NewManagedPrivateEndpointsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Managed Private Endpoints Client: %+v", err)
	}
	o.Configure(managedPrivateEndpointsClient.Client, o.Authorizers.Synapse)

	roleAssignmentsClient, err := roleassignments.NewRoleAssignmentsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Role Assignments Client: %+v", err)
	}
	o.Configure(roleAssignmentsClient.Client, o.Authorizers.Synapse)

	synapseRoleDefinitionsClient, err := synapseroledefinitions.NewSynapseRoleDefinitionsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Role Definitions Client: %+v", err)
	}
	o.Configure(synapseRoleDefinitionsClient.Client, o.Authorizers.Synapse)

	linkedServiceClient, err := linkedservices.NewLinkedServicesClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Linked Service Client: %+v", err)
	}
	o.Configure(linkedServiceClient.Client, o.Authorizers.Synapse)

	// TODO: migrate to go-azure-sdk
	firewallRuleClient, err := ipfirewallrules.NewIPFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Firewall Rules Client: %+v", err)
	}
	o.Configure(firewallRuleClient.Client, o.Authorizers.ResourceManager)

	integrationRuntimeAuthKeysClient := synapse.NewIntegrationRuntimeAuthKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationRuntimeAuthKeysClient.Client, o.ResourceManagerAuthorizer)

	integrationRuntimesClient := synapse.NewIntegrationRuntimesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationRuntimesClient.Client, o.ResourceManagerAuthorizer)

	keysClient, err := keys.NewKeysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Keys Client: %+v", err)
	}
	o.Configure(keysClient.Client, o.Authorizers.ResourceManager)

	privateLinkHubsClient, err := privatelinkhubs.NewPrivateLinkHubsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Private Link Hubs Client: %+v", err)
	}
	o.Configure(privateLinkHubsClient.Client, o.Authorizers.ResourceManager)

	// the service team hopes to rename it to sparkPool, so rename the sdk here
	sparkPoolClient, err := bigdatapools.NewBigDataPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Spark Pool Client: %+v", err)
	}
	o.Configure(sparkPoolClient.Client, o.Authorizers.ResourceManager)

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

	workspaceAzureADOnlyAuthenticationsClient := synapse.NewAzureADOnlyAuthenticationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceAzureADOnlyAuthenticationsClient.Client, o.ResourceManagerAuthorizer)

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
		// Resource Manager
		WorkspacesClient: workspacesClient,

		// Data Plane
		ManagedPrivateEndpointsClient: managedPrivateEndpointsClient,

		// TODO: Migrate to go-azure-sdk
		FirewallRulesClient:                               firewallRuleClient,
		IntegrationRuntimeAuthKeysClient:                  &integrationRuntimeAuthKeysClient,
		IntegrationRuntimesClient:                         &integrationRuntimesClient,
		KeysClient:                                        keysClient,
		PrivateLinkHubsClient:                             privateLinkHubsClient,
		SparkPoolClient:                                   sparkPoolClient,
		SqlPoolClient:                                     &sqlPoolClient,
		SqlPoolExtendedBlobAuditingPoliciesClient:         &sqlPoolExtendedBlobAuditingPoliciesClient,
		SqlPoolGeoBackupPoliciesClient:                    &sqlPoolGeoBackupPoliciesClient,
		SqlPoolSecurityAlertPolicyClient:                  &sqlPoolSecurityAlertPolicyClient,
		SqlPoolTransparentDataEncryptionClient:            &sqlPoolTransparentDataEncryptionClient,
		SqlPoolVulnerabilityAssessmentsClient:             &sqlPoolVulnerabilityAssessmentsClient,
		SQLPoolVulnerabilityAssessmentRuleBaselinesClient: &sqlPoolVulnerabilityAssessmentRuleBaselinesClient,
		SQLPoolWorkloadClassifierClient:                   &sqlPoolWorkloadClassifierClient,
		SQLPoolWorkloadGroupClient:                        &sqlPoolWorkloadGroupClient,
		WorkspaceAzureADOnlyAuthenticationsClient:         &workspaceAzureADOnlyAuthenticationsClient,
		WorkspaceAadAdminsClient:                          &workspaceAadAdminsClient,
		WorkspaceClient:                                   &workspaceClient,
		WorkspaceExtendedBlobAuditingPoliciesClient:       &workspaceExtendedBlobAuditingPoliciesClient,
		WorkspaceManagedIdentitySQLControlSettingsClient:  &workspaceManagedIdentitySQLControlSettingsClient,
		WorkspaceSecurityAlertPolicyClient:                &workspaceSecurityAlertPolicyClient,
		WorkspaceSQLAadAdminsClient:                       &workspaceSQLAadAdminsClient,
		WorkspaceVulnerabilityAssessmentsClient:           &workspaceVulnerabilityAssessmentsClient,

		// Data Plane (go-azure-sdk)
		roleAssignmentsClient:        roleAssignmentsClient,
		synapseRoleDefinitionsClient: synapseRoleDefinitionsClient,
		linkedServiceClient:          linkedServiceClient,
		synapseEnabled:               o.Authorizers.Synapse != nil,
	}, nil
}

func (client Client) RoleDefinitionsClient(workspaceName, synapseEndpointSuffix string) (*synapseroledefinitions.SynapseRoleDefinitionsClient, error) {
	if !client.synapseEnabled {
		return nil, errors.New("'Synapse' is not supported in this Azure Environment")
	}
	return client.synapseRoleDefinitionsClient.Clone(buildEndpoint(workspaceName, synapseEndpointSuffix)), nil
}

func (client Client) RoleAssignmentsClient(workspaceName, synapseEndpointSuffix string) (*roleassignments.RoleAssignmentsClient, error) {
	if !client.synapseEnabled {
		return nil, errors.New("'Synapse' is not supported in this Azure Environment")
	}
	return client.roleAssignmentsClient.Clone(buildEndpoint(workspaceName, synapseEndpointSuffix)), nil
}

func (client Client) LinkedServiceClient(workspaceName, synapseEndpointSuffix string) (*linkedservices.LinkedServicesClient, error) {
	if !client.synapseEnabled {
		return nil, errors.New("'Synapse' is not supported in this Azure Environment")
	}
	return client.linkedServiceClient.Clone(buildEndpoint(workspaceName, synapseEndpointSuffix)), nil
}

func buildEndpoint(workspaceName string, synapseEndpointSuffix string) string {
	return fmt.Sprintf("https://%s.%s", workspaceName, synapseEndpointSuffix)
}
