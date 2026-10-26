// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package preflight_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/appserviceplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight/testdata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// Built-in policy definitions
const appServiceAppsShouldUseASKUThatSupportsPrivateLink = "/providers/Microsoft.Authorization/policyDefinitions/546fe8d2-368d-4029-a418-6af48a7f61e5"
const allowedLocations = "/providers/Microsoft.Authorization/policyDefinitions/e56962a6-4747-49cd-b67b-bf8b01975c4c"

func TestAccPreflight_validateResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skipf("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	p := provider.AzureProvider()
	d := p.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d.HasError() {
		t.Fatalf("diags had error!")
	}
	client := p.Meta().(*clients.Client)
	skuPolicyResourceGroup := fmt.Sprintf("acctestRG-preflight-%d", time.Now().UnixNano())

	cases := []struct {
		Name          string
		Location      string
		ID            resourceids.ResourceId
		APIVersion    string
		Properties    any
		SKU           *any
		ExpectError   bool
		ErrorContains string
		BeforeTest    func(*testing.T)
	}{
		{
			Name:        "valid_webapp",
			Location:    "northeurope",
			ID:          pointer.To(commonids.NewAppServiceID(client.Account.SubscriptionId, "testResourceGroup", "testSiteName")),
			APIVersion:  "2024-11-01",
			Properties:  testdata.WebAppExample(client.Account.SubscriptionId),
			ExpectError: false,
		},
		{
			Name:        "valid_vnet",
			Location:    "northeurope",
			ID:          pointer.To(commonids.NewVirtualNetworkID(client.Account.SubscriptionId, "testResourceGroup", "testVirtualValid")),
			APIVersion:  "2024-01-01",
			Properties:  testdata.VirtualNetworkExample(),
			ExpectError: false,
		},
		{
			Name:        "invalid_vnet_missing_required_property",
			Location:    "northeurope",
			ID:          pointer.To(commonids.NewVirtualNetworkID(client.Account.SubscriptionId, "testResourceGroup", "testVirtualValid")),
			APIVersion:  "2024-01-01",
			Properties:  testdata.VirtualNetworkExampleMissingRequiredProperty(),
			ExpectError: true,
		},
		{
			Name:        "service_plan_p0v4_allowed_by_private_link_sku_policy",
			Location:    "westus2",
			ID:          pointer.To(commonids.NewAppServicePlanID(client.Account.SubscriptionId, skuPolicyResourceGroup, "preflight-p0v4")),
			APIVersion:  "2023-12-01",
			Properties:  testdata.AppServicePlanP0V4Example(),
			SKU:         pointer.To[any](&appserviceplans.SkuDescription{Name: pointer.To("P0V4")}),
			ExpectError: false,
			BeforeTest: func(t *testing.T) {
				groupID := setupResourceGroup(t, client, skuPolicyResourceGroup, "westus2")
				assignPolicy(t, client, groupID, appServiceAppsShouldUseASKUThatSupportsPrivateLink, fmt.Sprintf("preflight-sku-%d", time.Now().UnixNano()), nil)
			},
		},
		{
			Name:          "service_plan_f1_denied_by_private_link_sku_policy",
			Location:      "westus2",
			ID:            pointer.To(commonids.NewAppServicePlanID(client.Account.SubscriptionId, skuPolicyResourceGroup, "preflight-f1")),
			APIVersion:    "2023-12-01",
			Properties:    testdata.AppServicePlanF1Example(),
			SKU:           pointer.To[any](&appserviceplans.SkuDescription{Name: pointer.To("F1")}),
			ExpectError:   true,
			ErrorContains: appServiceAppsShouldUseASKUThatSupportsPrivateLink,
			BeforeTest: func(t *testing.T) {
				groupID := setupResourceGroup(t, client, skuPolicyResourceGroup, "westus2")
				assignPolicy(t, client, groupID, appServiceAppsShouldUseASKUThatSupportsPrivateLink, fmt.Sprintf("preflight-sku-%d", time.Now().UnixNano()), nil)
			},
		},
	}

	metadata := sdk.ResourceMetaData{
		Client: client,
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.BeforeTest != nil {
				tc.BeforeTest(t)
			}
			testCtx, testCancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer testCancel()
			request, err := preflight.NewValidationRequest(pointer.To(tc.Location), tc.ID, tc.APIVersion, tc.Properties, tc.SKU)
			if err != nil {
				t.Fatalf("building validation request: %s", err)
			}

			err = request.ValidateResource(testCtx, metadata)
			if err != nil && !tc.ExpectError {
				t.Fatalf("expected no error, got: %s", err)
			}
			if err == nil && tc.ExpectError {
				t.Fatalf("expected an error but got none")
			}
			if err != nil && tc.ErrorContains != "" && !strings.Contains(err.Error(), tc.ErrorContains) {
				t.Fatalf("expected error containing %q, got: %s", tc.ErrorContains, err)
			}
		})
	}
}

func setupResourceGroup(t *testing.T, client *clients.Client, resourceGroupName, location string) commonids.ResourceGroupId {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	groupID := commonids.NewResourceGroupID(client.Account.SubscriptionId, resourceGroupName)
	existing, err := client.Resource.ResourceGroupsClient.Get(ctx, groupID)
	if err == nil {
		return groupID
	}
	if !response.WasNotFound(existing.HttpResponse) {
		t.Fatalf("checking resource group: %s", err)
	}
	if _, err := client.Resource.ResourceGroupsClient.CreateOrUpdate(ctx, groupID, resourcegroups.ResourceGroup{Location: location}); err != nil {
		t.Fatalf("creating resource group: %s", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := client.Resource.ResourceGroupsClient.DeleteThenPoll(ctx, groupID, resourcegroups.DefaultDeleteOperationOptions()); err != nil {
			t.Errorf("deleting test resource group: %s", err)
		}
	})
	return groupID
}

func assignPolicy(t *testing.T, client *clients.Client, groupID commonids.ResourceGroupId, policyDefinitionID, assignmentName string, parameters map[string]policyassignments.ParameterValuesValue) policyassignments.ScopedPolicyAssignmentId {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	assignmentID := policyassignments.NewScopedPolicyAssignmentID(groupID.ID(), assignmentName)
	assignmentParameters := map[string]policyassignments.ParameterValuesValue{
		"effect": {Value: pointer.To[any]("Deny")},
	}
	for name, value := range parameters {
		assignmentParameters[name] = value
	}
	assignment := policyassignments.PolicyAssignment{
		Properties: &policyassignments.PolicyAssignmentProperties{
			PolicyDefinitionId: pointer.To(policyDefinitionID),
			Scope:              pointer.To(assignmentID.Scope),
			EnforcementMode:    pointer.To(policyassignments.EnforcementModeDefault),
			Parameters:         pointer.To(assignmentParameters),
		},
	}
	if _, err := client.Policy.AssignmentsClient.Create(ctx, assignmentID, assignment); err != nil {
		t.Fatalf("creating policy assignment: %s", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if resp, err := client.Policy.AssignmentsClient.Delete(ctx, assignmentID); err != nil && !response.WasNotFound(resp.HttpResponse) {
			t.Errorf("deleting policy assignment: %s", err)
		}
	})
	return assignmentID
}
