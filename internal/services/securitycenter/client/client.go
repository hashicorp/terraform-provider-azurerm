// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-12-01-preview/defenderforstorage"
	pricings_v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssessmentsClient                   *security.AssessmentsClient
	AssessmentsMetadataClient           *assessmentsmetadata.AssessmentsMetadataClient
	ContactsClient                      *security.ContactsClient
	DeviceSecurityGroupsClient          *security.DeviceSecurityGroupsClient
	IotSecuritySolutionClient           *security.IotSecuritySolutionClient
	PricingClient                       *pricings_v2023_01_01.PricingsClient
	WorkspaceClient                     *security.WorkspaceSettingsClient
	AdvancedThreatProtectionClient      *security.AdvancedThreatProtectionClient
	AutoProvisioningClient              *security.AutoProvisioningSettingsClient
	SettingClient                       *security.SettingsClient
	AutomationsClient                   *automations.AutomationsClient
	ServerVulnerabilityAssessmentClient *security.ServerVulnerabilityAssessmentClient
	DefenderForStorageClient            *defenderforstorage.DefenderForStorageClient
}

func NewClient(o *common.ClientOptions) *Client {
	ascLocation := "Global"

	AssessmentsClient := security.NewAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AssessmentsClient.Client, o.ResourceManagerAuthorizer)

	AssessmentsMetadataClient := assessmentsmetadata.NewAssessmentsMetadataClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AssessmentsMetadataClient.Client, o.ResourceManagerAuthorizer)

	ContactsClient := security.NewContactsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&ContactsClient.Client, o.ResourceManagerAuthorizer)

	DeviceSecurityGroupsClient := security.NewDeviceSecurityGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&DeviceSecurityGroupsClient.Client, o.ResourceManagerAuthorizer)

	IotSecuritySolutionClient := security.NewIotSecuritySolutionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&IotSecuritySolutionClient.Client, o.ResourceManagerAuthorizer)

	PricingClient := pricings_v2023_01_01.NewPricingsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&PricingClient.Client, o.ResourceManagerAuthorizer)

	WorkspaceClient := security.NewWorkspaceSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&WorkspaceClient.Client, o.ResourceManagerAuthorizer)

	AdvancedThreatProtectionClient := security.NewAdvancedThreatProtectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AdvancedThreatProtectionClient.Client, o.ResourceManagerAuthorizer)

	AutoProvisioningClient := security.NewAutoProvisioningSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AutoProvisioningClient.Client, o.ResourceManagerAuthorizer)

	SettingClient := security.NewSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&SettingClient.Client, o.ResourceManagerAuthorizer)

	AutomationsClient := automations.NewAutomationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AutomationsClient.Client, o.ResourceManagerAuthorizer)

	ServerVulnerabilityAssessmentClient := security.NewServerVulnerabilityAssessmentClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&ServerVulnerabilityAssessmentClient.Client, o.ResourceManagerAuthorizer)

	DefenderForStorageClient := defenderforstorage.NewDefenderForStorageClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DefenderForStorageClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AssessmentsClient:                   &AssessmentsClient,
		AssessmentsMetadataClient:           &AssessmentsMetadataClient,
		ContactsClient:                      &ContactsClient,
		DeviceSecurityGroupsClient:          &DeviceSecurityGroupsClient,
		IotSecuritySolutionClient:           &IotSecuritySolutionClient,
		PricingClient:                       &PricingClient,
		WorkspaceClient:                     &WorkspaceClient,
		AdvancedThreatProtectionClient:      &AdvancedThreatProtectionClient,
		AutoProvisioningClient:              &AutoProvisioningClient,
		SettingClient:                       &SettingClient,
		AutomationsClient:                   &AutomationsClient,
		ServerVulnerabilityAssessmentClient: &ServerVulnerabilityAssessmentClient,
		DefenderForStorageClient:            &DefenderForStorageClient,
	}
}
