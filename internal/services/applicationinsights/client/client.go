// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	analyticsitems "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
	apikeys "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
	billing "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentfeaturesandpricingapis"
	smartdetection "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	workbooktemplates "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	workbooks "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	webtests "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AnalyticsItemsClient     *analyticsitems.AnalyticsItemsAPIsClient
	APIKeysClient            *apikeys.ComponentApiKeysAPIsClient
	ComponentsClient         *components.ComponentsAPIsClient
	WebTestsClient           *webtests.WebTestsAPIsClient
	StandardWebTestsClient   *webtests.WebTestsAPIsClient
	BillingClient            *billing.ComponentFeaturesAndPricingAPIsClient
	SmartDetectionRuleClient *smartdetection.ComponentProactiveDetectionAPIsClient
	WorkbookClient           *workbooks.WorkbooksAPIsClient
	WorkbookTemplateClient   *workbooktemplates.WorkbookTemplatesAPIsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	analyticsItemsClient, err := analyticsitems.NewAnalyticsItemsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AnalyticsItems client: %+v", err)
	}
	o.Configure(analyticsItemsClient.Client, o.Authorizers.ResourceManager)

	apiKeysClient, err := apikeys.NewComponentApiKeysAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ApiKeys client: %+v", err)
	}
	o.Configure(apiKeysClient.Client, o.Authorizers.ResourceManager)

	componentsClient, err := components.NewComponentsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Components client: %+v", err)
	}
	o.Configure(componentsClient.Client, o.Authorizers.ResourceManager)

	webTestsClient, err := webtests.NewWebTestsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WebTests client: %+v", err)
	}
	o.Configure(webTestsClient.Client, o.Authorizers.ResourceManager)

	standardWebTestsClient, err := webtests.NewWebTestsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StandardWebTests client: %+v", err)
	}
	o.Configure(standardWebTestsClient.Client, o.Authorizers.ResourceManager)

	billingClient, err := billing.NewComponentFeaturesAndPricingAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Billing client: %+v", err)
	}
	o.Configure(billingClient.Client, o.Authorizers.ResourceManager)

	smartDetectionRuleClient, err := smartdetection.NewComponentProactiveDetectionAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SmartDetection client: %+v", err)
	}
	o.Configure(smartDetectionRuleClient.Client, o.Authorizers.ResourceManager)

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
