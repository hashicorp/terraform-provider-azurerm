// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	ServicePlanTypeConsumption = "consumption"
	ServicePlanTypeElastic     = "elastic"
	ServicePlanTypeIsolated    = "isolated"
	ServicePlanTypeAppPlan     = "app"
)

var appServicePlanSkus = []string{
	"B1", "B2", "B3", // basic
	"S1", "S2", "S3", // standard
	"P1v2", "P2v2", "P3v2", // Premium V2
	"P0v3", "P1v3", "P2v3", "P3v3", // Premium V3
	"P1mv3", "P2mv3", "P3mv3", "P4mv3", "P5mv3", // Premium V3 memory optimized
}

var freeSkus = []string{
	"F1",
}

var sharedSkus = []string{
	"D1",
	"SHARED",
}

var consumptionSkus = []string{
	"Y1",
}

var elasticSkus = []string{
	"EP1", "EP2", "EP3",
}

var isolatedSkus = []string{
	"I1", "I2", "I3", // Isolated V1 - ASEV2
	"I1v2", "I2v2", "I3v2", "I4v2", "I5v2", "I6v2", // Isolated v2 - ASEv3
}

var workflowSkus = []string{
	"WS1", "WS2", "WS3",
}

// AllKnownServicePlanSkus returns a list of all supported known SKU names
func AllKnownServicePlanSkus() []string {
	allSkus := make([]string, 0)
	allSkus = append(allSkus, appServicePlanSkus...)
	allSkus = append(allSkus, consumptionSkus...)
	allSkus = append(allSkus, elasticSkus...)
	allSkus = append(allSkus, freeSkus...)
	allSkus = append(allSkus, isolatedSkus...)
	allSkus = append(allSkus, sharedSkus...)
	allSkus = append(allSkus, workflowSkus...)

	return allSkus
}

func PlanIsConsumption(input *string) bool {
	if input == nil {
		return false
	}
	for _, v := range consumptionSkus {
		if strings.EqualFold(*input, v) {
			return true
		}
	}

	return false
}

func PlanIsElastic(input *string) bool {
	if input == nil {
		return false
	}
	for _, v := range elasticSkus {
		if strings.EqualFold(*input, v) {
			return true
		}
	}

	return false
}

func PlanIsIsolated(input *string) bool {
	if input == nil {
		return false
	}
	for _, v := range isolatedSkus {
		if strings.EqualFold(*input, v) {
			return true
		}
	}

	return false
}

func PlanIsAppPlan(input *string) bool {
	if input == nil {
		return false
	}
	for _, v := range appServicePlanSkus {
		if strings.EqualFold(*input, v) {
			return true
		}
	}

	return false
}

func PlanTypeFromSku(input string) string {
	if PlanIsConsumption(&input) {
		return ServicePlanTypeConsumption
	}

	if PlanIsElastic(&input) {
		return ServicePlanTypeElastic
	}

	if PlanIsIsolated(&input) {
		return ServicePlanTypeIsolated
	}

	if PlanIsAppPlan(&input) {
		return ServicePlanTypeAppPlan
	}

	return "unknown"
}

// ServicePlanInfoForApp returns the OS type and Service Plan SKU for a given App Service Resource
func ServicePlanInfoForApp(ctx context.Context, metadata sdk.ResourceMetaData, id interface{}) (osType *string, planSku *string, err error) {
	client := metadata.Client.AppService.WebAppsClient
	servicePlanClient := metadata.Client.AppService.ServicePlanClient
	var rg, siteName string

	switch appId := id.(type) {
	case parse.WebAppId:
		rg = appId.ResourceGroup
		siteName = appId.SiteName
	case parse.WebAppSlotId:
		rg = appId.ResourceGroup
		siteName = appId.SiteName
	case parse.FunctionAppId:
		rg = appId.ResourceGroup
		siteName = appId.SiteName
	case parse.FunctionAppSlotId:
		rg = appId.ResourceGroup
		siteName = appId.SiteName
	}

	site, err := client.Get(ctx, rg, siteName)
	if err != nil || site.SiteProperties == nil {
		return nil, nil, fmt.Errorf("reading %s: %+v", id, err)
	}
	props := site.SiteProperties
	if props.ServerFarmID == nil {
		return nil, nil, fmt.Errorf("determining Service Plan ID for %s: %+v", id, err)
	}
	servicePlanId, err := parse.ServicePlanID(*props.ServerFarmID)
	if err != nil {
		return nil, nil, err
	}

	sp, err := servicePlanClient.Get(ctx, servicePlanId.ResourceGroup, servicePlanId.ServerfarmName)
	if err != nil || sp.Kind == nil {
		return nil, nil, fmt.Errorf("reading Service Plan for %s: %+v", id, err)
	}

	osType = utils.String("windows")
	if strings.Contains(strings.ToLower(*sp.Kind), "linux") {
		osType = utils.String("linux")
	}

	planSku = utils.String("")
	if sku := sp.Sku; sku != nil {
		planSku = sku.Name
	}

	return osType, planSku, nil
}
