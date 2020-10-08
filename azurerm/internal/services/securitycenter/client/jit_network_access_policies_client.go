package client

import (
	"context"

	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
)

const fqdn = "github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"

type CustomJitNetworkAccessPoliciesClient struct {
	security.JitNetworkAccessPoliciesClient
}

func NewCustomJitNetworkAccessPoliciesClientWithBaseURI(baseURI string, subscriptionID string, ascLocation string) CustomJitNetworkAccessPoliciesClient {
	return CustomJitNetworkAccessPoliciesClient{security.NewJitNetworkAccessPoliciesClientWithBaseURI(baseURI, subscriptionID, ascLocation)}
}

func (client CustomJitNetworkAccessPoliciesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, jitNetworkAccessPolicyName string, ascLocation string, body security.JitNetworkAccessPolicy) (result security.JitNetworkAccessPolicy, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/JitNetworkAccessPoliciesClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: body,
			Constraints: []validation.Constraint{{Target: "body.JitNetworkAccessPolicyProperties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "body.JitNetworkAccessPolicyProperties.VirtualMachines", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("security.JitNetworkAccessPoliciesClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.RequestPreparer(autorest.AsPut(), ctx, resourceGroupName, jitNetworkAccessPolicyName, ascLocation, &body)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

func (client CustomJitNetworkAccessPoliciesClient) Delete(ctx context.Context, resourceGroupName string, jitNetworkAccessPolicyName string, ascLocation string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/JitNetworkAccessPoliciesClient.Delete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.JitNetworkAccessPoliciesClient", "Delete", err.Error())
	}

	req, err := client.RequestPreparer(autorest.AsDelete(), ctx, resourceGroupName, jitNetworkAccessPolicyName, ascLocation, nil)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "Delete", resp, "Failure responding to request")
	}

	return
}

func (client CustomJitNetworkAccessPoliciesClient) Get(ctx context.Context, resourceGroupName string, jitNetworkAccessPolicyName string, ascLocation string) (result security.JitNetworkAccessPolicy, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/JitNetworkAccessPoliciesClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.JitNetworkAccessPoliciesClient", "Get", err.Error())
	}

	req, err := client.RequestPreparer(autorest.AsGet(), ctx, resourceGroupName, jitNetworkAccessPolicyName, ascLocation, nil)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.JitNetworkAccessPoliciesClient", "Get", resp, "Failure responding to request")
	}

	return
}

func (client CustomJitNetworkAccessPoliciesClient) RequestPreparer(method autorest.PrepareDecorator, ctx context.Context, resourceGroupName string, jitNetworkAccessPolicyName string, ascLocation string, body *security.JitNetworkAccessPolicy) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"ascLocation":                autorest.Encode("path", ascLocation),
		"jitNetworkAccessPolicyName": autorest.Encode("path", jitNetworkAccessPolicyName),
		"resourceGroupName":          autorest.Encode("path", resourceGroupName),
		"subscriptionId":             autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		method,
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/locations/{ascLocation}/jitNetworkAccessPolicies/{jitNetworkAccessPolicyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if body != nil {
		preparer = autorest.DecoratePreparer(preparer, autorest.WithJSON(*body))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
