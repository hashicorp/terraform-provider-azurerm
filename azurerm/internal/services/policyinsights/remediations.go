package policyinsights

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/policyinsights/mgmt/2019-10-01/policyinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type managementGroupId struct {
	groupId string
}

type subscriptionId struct {
	subscriptionId string
}

func parseManagementGroupId(input string) (*managementGroupId, error) {
	// /providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000
	input = strings.Trim(input, "/")
	inputLower := strings.ToLower(input)

	const prefix = "providers/microsoft.management"
	if !strings.HasPrefix(inputLower, prefix) {
		return nil, fmt.Errorf("Expected input to start with %s", prefix)
	}

	segments := strings.Split(input, "/")
	if len(segments) != 4 {
		return nil, fmt.Errorf("Expected there to be 4 segments but got %d", len(segments))
	}

	id := managementGroupId{
		groupId: segments[3],
	}
	return &id, nil
}

func parseSubscriptionID(input string) (*subscriptionId, error) {
	// this is either:
	// /subscriptions/00000000-0000-0000-0000-000000000000
	subString := strings.Trim(input, "/")
	inputLower := strings.ToLower(subString)

	components := strings.Split(inputLower, "/")

	if len(components) == 0 {
		return nil, fmt.Errorf("Subscription Id is empty or not formatted correctly: %s", subString)
	}

	if len(components) != 2 {
		return nil, fmt.Errorf("Subscription Id should have 2 segments, got %d: %q", len(components), subString)
	}

	if components[0] != "subscriptions" {
		return nil, fmt.Errorf("%q is not legal subscription ID", input)
	}

	id := subscriptionId{
		subscriptionId: components[1],
	}
	return &id, nil
}

type ProvisioningType int

const (
	AtSubscription ProvisioningType = iota
	AtManagementGroup
	AtResourceGroup
	AtResource
)

type RemediationId struct {
	Name string
	RemediationScope
}

type RemediationScope struct {
	Type              ProvisioningType
	Scope             string
	SubscriptionId    *string
	ManagementGroupId *string
	ResourceGroup     *string
}

func ParseScope(scope string) (*RemediationScope, error) {
	scopeObj := RemediationScope{
		Scope: scope,
	}
	if subsId, err := parseSubscriptionID(scope); err == nil && subsId != nil {
		log.Printf("[INFO] Get Subscription ID from scope %s", scope)
		scopeObj.SubscriptionId = &subsId.subscriptionId
		scopeObj.Type = AtSubscription
	} else if mgmtId, err := parseManagementGroupId(scope); err == nil && mgmtId != nil {
		log.Printf("[INFO] Get Management Group ID from scope %s", scope)
		scopeObj.ManagementGroupId = &mgmtId.groupId
		scopeObj.Type = AtManagementGroup
	} else if resId, err := azure.ParseAzureResourceID(scope); err == nil && resId != nil {
		if len(resId.Path) == 0 {
			log.Printf("[INFO] Get Resource Group ID from scope %s", scope)
			scopeObj.SubscriptionId = &resId.SubscriptionID
			scopeObj.ResourceGroup = &resId.ResourceGroup
			scopeObj.Type = AtResourceGroup
		} else {
			log.Printf("[INFO] Get Resource ID from scope %s", scope)
			scopeObj.Type = AtResource
		}
	} else {
		return nil, fmt.Errorf("Cannot recognize scope %q as Subscription ID, Management Group ID, Resource Group ID, or Resource ID", scope)
	}
	return &scopeObj, nil
}

func ParseRemediationId(id string) (*RemediationId, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure ID: %+v", err)
	}

	path := idURL.Path
	path = strings.TrimSuffix(path, "/") // we need to keep the leading slash to keep the azure.ParseAzureResourceID working
	lowerPath := strings.ToLower(path)   // the id returned from service may not keep the casing

	const policyInsightProvider = "/providers/microsoft.policyinsights"
	index := strings.LastIndex(lowerPath, policyInsightProvider)
	if index < 0 {
		return nil, fmt.Errorf("Cannot find %s in path %s", policyInsightProvider, path)
	}

	scope := path[0:index]
	path = path[index+1:] // throw away the slash at beginning
	segments := strings.Split(path, "/")
	if len(segments) != 4 {
		return nil, fmt.Errorf("Expect id of remediation has 4 segment")
	}
	name := segments[3]

	scopeObj, err := ParseScope(scope)
	if err != nil {
		return nil, err
	}

	idObj := RemediationId{
		Name:             name,
		RemediationScope: *scopeObj,
	}

	return &idObj, nil
}

func remediationCreateUpdateAtResourceGroup(client *policyinsights.RemediationsClient, ctx context.Context, scopeObj *RemediationScope, d *schema.ResourceData) error {
	name := d.Get("name").(string)
	subscriptionId := *scopeObj.SubscriptionId
	resourceGroupName := *scopeObj.ResourceGroup

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetAtResourceGroup(ctx, subscriptionId, resourceGroupName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Policy Remediation %q (Resource Group Name %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_remediation", *existing.ID)
		}
	}

	filters := d.Get("location_filters").([]interface{})
	policyAssignmentID := d.Get("policy_assignment_id").(string)
	policyDefinitionReferenceID := d.Get("policy_definition_reference_id").(string)

	parameters := policyinsights.Remediation{
		RemediationProperties: &policyinsights.RemediationProperties{
			Filters:                     expandArmRemediationLocationFilters(filters),
			PolicyAssignmentID:          utils.String(policyAssignmentID),
			PolicyDefinitionReferenceID: utils.String(policyDefinitionReferenceID),
		},
	}

	if _, err := client.CreateOrUpdateAtResourceGroup(ctx, subscriptionId, resourceGroupName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Policy Remediation %q (Resource Group Name %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.GetAtResourceGroup(ctx, subscriptionId, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Policy Remediation %q (Resource Group Name %q): %+v", name, resourceGroupName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Policy Remediation %q (Resource Group Name %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return nil
}

func remediationCreateUpdateAtSubscription(client *policyinsights.RemediationsClient, ctx context.Context, scope *RemediationScope, d *schema.ResourceData) error {
	name := d.Get("name").(string)
	subscriptionId := *scope.SubscriptionId

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetAtSubscription(ctx, subscriptionId, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Policy Remediation %q (Subscription ID %q): %+v", name, subscriptionId, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_remediation", *existing.ID)
		}
	}

	filters := d.Get("location_filters").([]interface{})
	policyAssignmentID := d.Get("policy_assignment_id").(string)
	policyDefinitionReferenceID := d.Get("policy_definition_reference_id").(string)

	parameters := policyinsights.Remediation{
		RemediationProperties: &policyinsights.RemediationProperties{
			Filters:                     expandArmRemediationLocationFilters(filters),
			PolicyAssignmentID:          utils.String(policyAssignmentID),
			PolicyDefinitionReferenceID: utils.String(policyDefinitionReferenceID),
		},
	}

	if _, err := client.CreateOrUpdateAtSubscription(ctx, subscriptionId, name, parameters); err != nil {
		return fmt.Errorf("Error creating Policy Remediation %q (Subscription ID %q): %+v", name, subscriptionId, err)
	}

	resp, err := client.GetAtSubscription(ctx, subscriptionId, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Policy Remediation %q (Subscription ID %q): %+v", name, subscriptionId, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Policy Remediation %q (Subscription ID %q) ID", name, subscriptionId)
	}
	d.SetId(*resp.ID)

	return nil
}

func remediationCreateUpdateAtManagementGroup(client *policyinsights.RemediationsClient, ctx context.Context, scope *RemediationScope, d *schema.ResourceData) error {
	name := d.Get("name").(string)
	groupId := *scope.ManagementGroupId

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetAtManagementGroup(ctx, groupId, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Policy Remediation %q (Management Group %q): %+v", name, groupId, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_remediation", *existing.ID)
		}
	}

	filters := d.Get("location_filters").([]interface{})
	policyAssignmentID := d.Get("policy_assignment_id").(string)
	policyDefinitionReferenceID := d.Get("policy_definition_reference_id").(string)

	parameters := policyinsights.Remediation{
		RemediationProperties: &policyinsights.RemediationProperties{
			Filters:                     expandArmRemediationLocationFilters(filters),
			PolicyAssignmentID:          utils.String(policyAssignmentID),
			PolicyDefinitionReferenceID: utils.String(policyDefinitionReferenceID),
		},
	}

	if _, err := client.CreateOrUpdateAtManagementGroup(ctx, groupId, name, parameters); err != nil {
		return fmt.Errorf("Error creating Policy Remediation %q (Management Group %q): %+v", name, groupId, err)
	}

	resp, err := client.GetAtManagementGroup(ctx, groupId, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Policy Remediation %q (Management Group %q): %+v", name, groupId, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Policy Remediation %q (Management Group %q) ID", name, groupId)
	}
	d.SetId(*resp.ID)

	return nil
}

func remediationCreateUpdateAtResource(client *policyinsights.RemediationsClient, ctx context.Context, scope *RemediationScope, d *schema.ResourceData) error {
	name := d.Get("name").(string)
	resourceId := scope.Scope

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetAtResource(ctx, resourceId, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Policy Remediation %q (Resource %q): %+v", name, resourceId, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_remediation", *existing.ID)
		}
	}

	filters := d.Get("location_filters").([]interface{})
	policyAssignmentID := d.Get("policy_assignment_id").(string)
	policyDefinitionReferenceID := d.Get("policy_definition_reference_id").(string)

	parameters := policyinsights.Remediation{
		RemediationProperties: &policyinsights.RemediationProperties{
			Filters:                     expandArmRemediationLocationFilters(filters),
			PolicyAssignmentID:          utils.String(policyAssignmentID),
			PolicyDefinitionReferenceID: utils.String(policyDefinitionReferenceID),
		},
	}

	if _, err := client.CreateOrUpdateAtResource(ctx, scope.Scope, name, parameters); err != nil {
		return fmt.Errorf("Error creating Policy Remediation %q (Resource %q): %+v", name, resourceId, err)
	}

	resp, err := client.GetAtResource(ctx, scope.Scope, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Policy Remediation %q (Resource %q): %+v", name, resourceId, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Policy Remediation %q (Resource %q) ID", name, resourceId)
	}
	d.SetId(*resp.ID)

	return nil
}
