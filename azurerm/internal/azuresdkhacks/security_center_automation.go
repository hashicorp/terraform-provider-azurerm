package azuresdkhacks

//
// Patched versions of: services/preview/security/mgmt/v1.0/security
// With Client, Preparer and main Automation struct
//
// Fixes for this issue: https://github.com/Azure/azure-sdk-for-go/issues/12634
//

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
)

// Automation wraps security.Automation
type Automation struct {
	security.Automation
}

// AutomationsClient wraps security.AutomationsClient
type AutomationsClient struct {
	security.AutomationsClient
}

// NewAutomationsClientWithBaseURI constructs a new patched version of AutomationsClient
func NewAutomationsClientWithBaseURI(baseURI string, azureSubID string, location string) AutomationsClient {
	return AutomationsClient{
		security.NewAutomationsClientWithBaseURI(baseURI, azureSubID, location),
	}
}

// MarshalJSON fixes the bug in the standard MarshalJSON func
func (a Automation) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if a.AutomationProperties != nil {
		objectMap["properties"] = a.AutomationProperties
	}
	if a.Kind != nil {
		objectMap["kind"] = a.Kind
	}
	if a.Etag != nil {
		objectMap["etag"] = a.Etag
	}
	if a.Tags != nil {
		objectMap["tags"] = a.Tags
	}
	// !FIXED! This is missing from the function in the SDK
	if a.Location != nil {
		objectMap["location"] = a.Location
	}
	return json.Marshal(objectMap)
}

// CreateOrUpdate with our patched client
func (client AutomationsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, automationName string, automation Automation) (result security.Automation, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.AutomationsClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, automationName, automation)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.AutomationsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.AutomationsClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.AutomationsClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer is used by the patched CreateOrUpdate
func (client AutomationsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, automationName string, automation Automation) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"automationName":    autorest.Encode("path", automationName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-01-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/automations/{automationName}", pathParameters),
		// This will result in the patched version of MarshallJSON being called
		autorest.WithJSON(automation),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
