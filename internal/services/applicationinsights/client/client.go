// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2020-02-02/insights" // nolint: staticcheck
	workbooktemplates "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	workbooks "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	webtests "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/azuresdkhacks"
)

type Client struct {
	AnalyticsItemsClient     *insights.AnalyticsItemsClient
	APIKeysClient            *insights.APIKeysClient
	ComponentsClient         *insights.ComponentsClient
	WebTestsClient           *azuresdkhacks.WebTestsClient
	StandardWebTestsClient   *webtests.WebTestsAPIsClient
	BillingClient            *insights.ComponentCurrentBillingFeaturesClient
	SmartDetectionRuleClient *insights.ProactiveDetectionConfigurationsClient
	WorkbookClient           *workbooks.WorkbooksAPIsClient
	WorkbookTemplateClient   *workbooktemplates.WorkbookTemplatesAPIsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	analyticsItemsClient := insights.NewAnalyticsItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&analyticsItemsClient.Client, o.ResourceManagerAuthorizer)

	apiKeysClient := insights.NewAPIKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiKeysClient.Client, o.ResourceManagerAuthorizer)

	componentsClient := insights.NewComponentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&componentsClient.Client, o.ResourceManagerAuthorizer)

	webTestsClient := insights.NewWebTestsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webTestsClient.Client, o.ResourceManagerAuthorizer)
	webTestsWorkaroundClient := azuresdkhacks.NewWebTestsClient(webTestsClient)

	standardWebTestsClient, err := webtests.NewWebTestsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StandardWebTests client: %+v", err)
	}
	o.Configure(standardWebTestsClient.Client, o.Authorizers.ResourceManager)

	billingClient := insights.NewComponentCurrentBillingFeaturesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&billingClient.Client, o.ResourceManagerAuthorizer)

	smartDetectionRuleClient := insights.NewProactiveDetectionConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&smartDetectionRuleClient.Client, o.ResourceManagerAuthorizer)

	workbookClient, err := workbooks.NewWorkbooksAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workbook client: %+v", err)
	}
	o.Configure(workbookClient.Client, o.Authorizers.ResourceManager)

	workbookTemplateClient, err := workbooktemplates.NewWorkbookTemplatesAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WorkbookTemplate client: %+v", err)
	}
	o.Configure(workbookTemplateClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AnalyticsItemsClient:     &analyticsItemsClient,
		APIKeysClient:            &apiKeysClient,
		ComponentsClient:         &componentsClient,
		WebTestsClient:           &webTestsWorkaroundClient,
		BillingClient:            &billingClient,
		SmartDetectionRuleClient: &smartDetectionRuleClient,
		WorkbookClient:           workbookClient,
		WorkbookTemplateClient:   workbookTemplateClient,
		StandardWebTestsClient:   standardWebTestsClient,
	}, nil
}
