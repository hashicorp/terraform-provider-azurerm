// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ManagedDatabasesClient                           *sql.ManagedDatabasesClient
	ManagedInstancesClient                           *sql.ManagedInstancesClient
	ManagedInstancesLongTermRetentionPoliciesClient  *sql.ManagedInstanceLongTermRetentionPoliciesClient
	ManagedInstancesShortTermRetentionPoliciesClient *sql.ManagedBackupShortTermRetentionPoliciesClient
	ManagedInstanceVulnerabilityAssessmentsClient    *sql.ManagedInstanceVulnerabilityAssessmentsClient
	ManagedInstanceServerSecurityAlertPoliciesClient *sql.ManagedServerSecurityAlertPoliciesClient
	ManagedInstanceAdministratorsClient              *sql.ManagedInstanceAdministratorsClient
	ManagedInstanceAzureADOnlyAuthenticationsClient  *sql.ManagedInstanceAzureADOnlyAuthenticationsClient
	ManagedInstanceEncryptionProtectorClient         *sql.ManagedInstanceEncryptionProtectorsClient
	ManagedInstanceFailoverGroupsClient              *sql.InstanceFailoverGroupsClient
	ManagedInstanceKeysClient                        *sql.ManagedInstanceKeysClient
}

func NewClient(o *common.ClientOptions) *Client {

	managedDatabasesClient := sql.NewManagedDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedDatabasesClient.Client, o.ResourceManagerAuthorizer)

	managedInstancesClient := sql.NewManagedInstancesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstancesClient.Client, o.ResourceManagerAuthorizer)

	managedInstancesLongTermRetentionPoliciesClient := sql.NewManagedInstanceLongTermRetentionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstancesLongTermRetentionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	managedInstancesShortTermRetentionPoliciesClient := sql.NewManagedBackupShortTermRetentionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstancesShortTermRetentionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	managedInstancesAdministratorsClient := sql.NewManagedInstanceAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstancesAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	managedInstanceAzureADOnlyAuthenticationsClient := sql.NewManagedInstanceAzureADOnlyAuthenticationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstanceAzureADOnlyAuthenticationsClient.Client, o.ResourceManagerAuthorizer)

	managedInstanceEncryptionProtectorsClient := sql.NewManagedInstanceEncryptionProtectorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstanceEncryptionProtectorsClient.Client, o.ResourceManagerAuthorizer)

	managedInstanceFailoverGroupsClient := sql.NewInstanceFailoverGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstanceFailoverGroupsClient.Client, o.ResourceManagerAuthorizer)

	managedInstanceKeysClient := sql.NewManagedInstanceKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstanceKeysClient.Client, o.ResourceManagerAuthorizer)

	managedInstanceVulnerabilityAssessmentsClient := sql.NewManagedInstanceVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstanceVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	managedInstanceServerSecurityAlertPoliciesClient := sql.NewManagedServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstanceServerSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagedDatabasesClient:                           &managedDatabasesClient,
		ManagedInstanceAdministratorsClient:              &managedInstancesAdministratorsClient,
		ManagedInstanceAzureADOnlyAuthenticationsClient:  &managedInstanceAzureADOnlyAuthenticationsClient,
		ManagedInstanceEncryptionProtectorClient:         &managedInstanceEncryptionProtectorsClient,
		ManagedInstanceFailoverGroupsClient:              &managedInstanceFailoverGroupsClient,
		ManagedInstanceKeysClient:                        &managedInstanceKeysClient,
		ManagedInstancesLongTermRetentionPoliciesClient:  &managedInstancesLongTermRetentionPoliciesClient,
		ManagedInstanceServerSecurityAlertPoliciesClient: &managedInstanceServerSecurityAlertPoliciesClient,
		ManagedInstancesShortTermRetentionPoliciesClient: &managedInstancesShortTermRetentionPoliciesClient,
		ManagedInstanceVulnerabilityAssessmentsClient:    &managedInstanceVulnerabilityAssessmentsClient,
		ManagedInstancesClient:                           &managedInstancesClient,
	}
}
