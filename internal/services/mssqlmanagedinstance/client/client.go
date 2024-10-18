// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/instancefailovergroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedbackupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/manageddatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstanceadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstanceazureadonlyauthentications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstanceencryptionprotectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstancekeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstancelongtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstancevulnerabilityassessments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedserversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ManagedDatabasesClient                           *manageddatabases.ManagedDatabasesClient
	ManagedInstancesClient                           *managedinstances.ManagedInstancesClient
	ManagedInstancesLongTermRetentionPoliciesClient  *managedinstancelongtermretentionpolicies.ManagedInstanceLongTermRetentionPoliciesClient
	ManagedInstancesShortTermRetentionPoliciesClient *managedbackupshorttermretentionpolicies.ManagedBackupShortTermRetentionPoliciesClient
	ManagedInstanceVulnerabilityAssessmentsClient    *managedinstancevulnerabilityassessments.ManagedInstanceVulnerabilityAssessmentsClient
	ManagedInstanceServerSecurityAlertPoliciesClient *managedserversecurityalertpolicies.ManagedServerSecurityAlertPoliciesClient
	ManagedInstanceAdministratorsClient              *managedinstanceadministrators.ManagedInstanceAdministratorsClient
	ManagedInstanceAzureADOnlyAuthenticationsClient  *managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthenticationsClient
	ManagedInstanceEncryptionProtectorClient         *managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtectorsClient
	ManagedInstanceFailoverGroupsClient              *instancefailovergroups.InstanceFailoverGroupsClient
	ManagedInstanceKeysClient                        *managedinstancekeys.ManagedInstanceKeysClient

	options *common.ClientOptions
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	managedDatabasesClient, err := manageddatabases.NewManagedDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Databases Client: %+v", err)
	}
	o.Configure(managedDatabasesClient.Client, o.Authorizers.ResourceManager)

	managedInstancesClient, err := managedinstances.NewManagedInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Client: %+v", err)
	}
	o.Configure(managedInstancesClient.Client, o.Authorizers.ResourceManager)

	managedInstancesLongTermRetentionPoliciesClient, err := managedinstancelongtermretentionpolicies.NewManagedInstanceLongTermRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Long Term Retention Policies Client: %+v", err)
	}
	o.Configure(managedInstancesLongTermRetentionPoliciesClient.Client, o.Authorizers.ResourceManager)

	managedInstancesShortTermRetentionPoliciesClient, err := managedbackupshorttermretentionpolicies.NewManagedBackupShortTermRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Short Term Retention Policies Client: %+v", err)
	}
	o.Configure(managedInstancesShortTermRetentionPoliciesClient.Client, o.Authorizers.ResourceManager)

	managedInstancesAdministratorsClient, err := managedinstanceadministrators.NewManagedInstanceAdministratorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Administrators Client: %+v", err)
	}
	o.Configure(managedInstancesAdministratorsClient.Client, o.Authorizers.ResourceManager)

	managedInstanceAzureADOnlyAuthenticationsClient, err := managedinstanceazureadonlyauthentications.NewManagedInstanceAzureADOnlyAuthenticationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Azure ADOnly Authentications Client: %+v", err)
	}
	o.Configure(managedInstanceAzureADOnlyAuthenticationsClient.Client, o.Authorizers.ResourceManager)

	managedInstanceEncryptionProtectorsClient, err := managedinstanceencryptionprotectors.NewManagedInstanceEncryptionProtectorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Encryption Protectors Client: %+v", err)
	}
	o.Configure(managedInstanceEncryptionProtectorsClient.Client, o.Authorizers.ResourceManager)

	managedInstanceFailoverGroupsClient, err := instancefailovergroups.NewInstanceFailoverGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Failover Groups Client: %+v", err)
	}
	o.Configure(managedInstanceFailoverGroupsClient.Client, o.Authorizers.ResourceManager)

	managedInstanceKeysClient, err := managedinstancekeys.NewManagedInstanceKeysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Keys Client: %+v", err)
	}
	o.Configure(managedInstanceKeysClient.Client, o.Authorizers.ResourceManager)

	managedInstanceVulnerabilityAssessmentsClient, err := managedinstancevulnerabilityassessments.NewManagedInstanceVulnerabilityAssessmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Vulnerability Assessments Client: %+v", err)
	}
	o.Configure(managedInstanceVulnerabilityAssessmentsClient.Client, o.Authorizers.ResourceManager)

	managedInstanceServerSecurityAlertPoliciesClient, err := managedserversecurityalertpolicies.NewManagedServerSecurityAlertPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Instance Server Security Alert Policies Client: %+v", err)
	}
	o.Configure(managedInstanceServerSecurityAlertPoliciesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ManagedDatabasesClient:                           managedDatabasesClient,
		ManagedInstanceAdministratorsClient:              managedInstancesAdministratorsClient,
		ManagedInstanceAzureADOnlyAuthenticationsClient:  managedInstanceAzureADOnlyAuthenticationsClient,
		ManagedInstanceEncryptionProtectorClient:         managedInstanceEncryptionProtectorsClient,
		ManagedInstanceFailoverGroupsClient:              managedInstanceFailoverGroupsClient,
		ManagedInstanceKeysClient:                        managedInstanceKeysClient,
		ManagedInstancesLongTermRetentionPoliciesClient:  managedInstancesLongTermRetentionPoliciesClient,
		ManagedInstanceServerSecurityAlertPoliciesClient: managedInstanceServerSecurityAlertPoliciesClient,
		ManagedInstancesShortTermRetentionPoliciesClient: managedInstancesShortTermRetentionPoliciesClient,
		ManagedInstanceVulnerabilityAssessmentsClient:    managedInstanceVulnerabilityAssessmentsClient,
		ManagedInstancesClient:                           managedInstancesClient,

		options: o,
	}, nil
}

func (c Client) ManagedInstancesClientForSubscription(subscriptionID string) *sql.ManagedInstancesClient {
	// TODO: this method can be removed once this is moved to using `hashicorp/go-azure-sdk`
	managedInstancesClient := sql.NewManagedInstancesClientWithBaseURI(c.options.ResourceManagerEndpoint, subscriptionID)
	c.options.ConfigureClient(&managedInstancesClient.Client, c.options.ResourceManagerAuthorizer)
	return &managedInstancesClient
}
