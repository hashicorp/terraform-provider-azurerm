// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentfeaturesandpricingapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AnalyticsItemsClient     *analyticsitemsapis.AnalyticsItemsAPIsClient
	APIKeysClient            *componentapikeysapis.ComponentApiKeysAPIsClient
	ComponentsClient         *componentsapis.ComponentsAPIsClient
	WebTestsClient           *webtestsapis.WebTestsAPIsClient
	StandardWebTestsClient   *webtestsapis.WebTestsAPIsClient
	BillingClient            *componentfeaturesandpricingapis.ComponentFeaturesAndPricingAPIsClient
	SmartDetectionRuleClient *componentproactivedetectionapis.ComponentProactiveDetectionAPIsClient
	WorkbookClient           *workbooksapis.WorkbooksAPIsClient
	WorkbookTemplateClient   *workbooktemplatesapis.WorkbookTemplatesAPIsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	analyticsItemsClient, err := analyticsitemsapis.NewAnalyticsItemsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AnalyticsItems client: %+v", err)
	}
	o.Configure(analyticsItemsClient.Client, o.Authorizers.ResourceManager)

	apiKeysClient, err := componentapikeysapis.NewComponentApiKeysAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ApiKeys client: %+v", err)
	}
	o.Configure(apiKeysClient.Client, o.Authorizers.ResourceManager)

	componentsClient, err := componentsapis.NewComponentsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Components client: %+v", err)
	}
	o.Configure(componentsClient.Client, o.Authorizers.ResourceManager)

	webTestsClient, err := webtestsapis.NewWebTestsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WebTests client: %+v", err)
	}
	o.Configure(webTestsClient.Client, o.Authorizers.ResourceManager)

	standardWebTestsClient, err := webtestsapis.NewWebTestsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StandardWebTests client: %+v", err)
	}
	o.Configure(standardWebTestsClient.Client, o.Authorizers.ResourceManager)

	billingClient, err := componentfeaturesandpricingapis.NewComponentFeaturesAndPricingAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Billing client: %+v", err)
	}
	o.Configure(billingClient.Client, o.Authorizers.ResourceManager)

	smartDetectionRuleClient, err := componentproactivedetectionapis.NewComponentProactiveDetectionAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SmartDetection client: %+v", err)
	}
	o.Configure(smartDetectionRuleClient.Client, o.Authorizers.ResourceManager)

	workbookClient, err := workbooksapis.NewWorkbooksAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workbook client: %+v", err)
	}
	o.Configure(workbookClient.Client, o.Authorizers.ResourceManager)

	workbookTemplateClient, err := workbooktemplatesapis.NewWorkbookTemplatesAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WorkbookTemplate client: %+v", err)
	}
	o.Configure(workbookTemplateClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AnalyticsItemsClient:     analyticsItemsClient,
		APIKeysClient:            apiKeysClient,
		ComponentsClient:         componentsClient,
		WebTestsClient:           webTestsClient,
		BillingClient:            billingClient,
		SmartDetectionRuleClient: smartDetectionRuleClient,
		WorkbookClient:           workbookClient,
		WorkbookTemplateClient:   workbookTemplateClient,
		StandardWebTestsClient:   standardWebTestsClient,
	}, nil
}
