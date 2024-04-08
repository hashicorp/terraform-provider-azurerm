// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-12-01-preview/defenderforstorage"
	pricings_v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-05-01/servervulnerabilityassessmentssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssessmentsClient                          *security.AssessmentsClient
	AssessmentsMetadataClient                  *assessmentsmetadata.AssessmentsMetadataClient
	ContactsClient                             *security.ContactsClient
	DeviceSecurityGroupsClient                 *security.DeviceSecurityGroupsClient
	IotSecuritySolutionClient                  *security.IotSecuritySolutionClient
	PricingClient                              *pricings_v2023_01_01.PricingsClient
	WorkspaceClient                            *security.WorkspaceSettingsClient
	AdvancedThreatProtectionClient             *security.AdvancedThreatProtectionClient
	AutoProvisioningClient                     *security.AutoProvisioningSettingsClient
	SettingClient                              *settings.SettingsClient
	AutomationsClient                          *automations.AutomationsClient
	ServerVulnerabilityAssessmentClient        *security.ServerVulnerabilityAssessmentClient
	ServerVulnerabilityAssessmentSettingClient *servervulnerabilityassessmentssettings.ServerVulnerabilityAssessmentsSettingsClient
	DefenderForStorageClient                   *defenderforstorage.DefenderForStorageClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	ascLocation := "Global"

	AssessmentsClient := security.NewAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AssessmentsClient.Client, o.ResourceManagerAuthorizer)

	AssessmentsMetadataClient, err := assessmentsmetadata.NewAssessmentsMetadataClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Assessments Metadata client : %+v", err)
	}
	o.Configure(AssessmentsMetadataClient.Client, o.Authorizers.ResourceManager)

	ContactsClient := security.NewContactsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&ContactsClient.Client, o.ResourceManagerAuthorizer)

	DeviceSecurityGroupsClient := security.NewDeviceSecurityGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&DeviceSecurityGroupsClient.Client, o.ResourceManagerAuthorizer)

	IotSecuritySolutionClient := security.NewIotSecuritySolutionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&IotSecuritySolutionClient.Client, o.ResourceManagerAuthorizer)

	PricingClient, err := pricings_v2023_01_01.NewPricingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Pricing client : %+v", err)
	}
	o.Configure(PricingClient.Client, o.Authorizers.ResourceManager)

	WorkspaceClient := security.NewWorkspaceSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&WorkspaceClient.Client, o.ResourceManagerAuthorizer)

	AdvancedThreatProtectionClient := security.NewAdvancedThreatProtectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AdvancedThreatProtectionClient.Client, o.ResourceManagerAuthorizer)

	AutoProvisioningClient := security.NewAutoProvisioningSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AutoProvisioningClient.Client, o.ResourceManagerAuthorizer)

	SettingClient, err := settings.NewSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Setting client : %+v", err)
	}
	o.Configure(SettingClient.Client, o.Authorizers.ResourceManager)

	AutomationsClient, err := automations.NewAutomationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Automations client : %+v", err)
	}
	o.Configure(AutomationsClient.Client, o.Authorizers.ResourceManager)

	ServerVulnerabilityAssessmentClient := security.NewServerVulnerabilityAssessmentClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&ServerVulnerabilityAssessmentClient.Client, o.ResourceManagerAuthorizer)

	ServerVulnerabilityAssessmentSettingClient, err := servervulnerabilityassessmentssettings.NewServerVulnerabilityAssessmentsSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Vulnerability Assessment Setting client : %+v", err)
	}
	o.Configure(ServerVulnerabilityAssessmentSettingClient.Client, o.Authorizers.ResourceManager)

	DefenderForStorageClient, err := defenderforstorage.NewDefenderForStorageClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Defender For Storage client : %+v", err)
	}
	o.Configure(DefenderForStorageClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AssessmentsClient:                          &AssessmentsClient,
		AssessmentsMetadataClient:                  AssessmentsMetadataClient,
		ContactsClient:                             &ContactsClient,
		DeviceSecurityGroupsClient:                 &DeviceSecurityGroupsClient,
		IotSecuritySolutionClient:                  &IotSecuritySolutionClient,
		PricingClient:                              PricingClient,
		WorkspaceClient:                            &WorkspaceClient,
		AdvancedThreatProtectionClient:             &AdvancedThreatProtectionClient,
		AutoProvisioningClient:                     &AutoProvisioningClient,
		SettingClient:                              SettingClient,
		AutomationsClient:                          AutomationsClient,
		ServerVulnerabilityAssessmentClient:        &ServerVulnerabilityAssessmentClient,
		ServerVulnerabilityAssessmentSettingClient: ServerVulnerabilityAssessmentSettingClient,
		DefenderForStorageClient:                   DefenderForStorageClient,
	}, nil
}
