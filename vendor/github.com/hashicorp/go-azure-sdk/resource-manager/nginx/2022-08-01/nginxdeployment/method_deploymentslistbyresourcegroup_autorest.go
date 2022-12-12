package nginxdeployment

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]NginxDeployment

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DeploymentsListByResourceGroupOperationResponse, error)
}

type DeploymentsListByResourceGroupCompleteResult struct {
	Items []NginxDeployment
}

func (r DeploymentsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DeploymentsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp DeploymentsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DeploymentsListByResourceGroup ...
func (c NginxDeploymentClient) DeploymentsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp DeploymentsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForDeploymentsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDeploymentsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForDeploymentsListByResourceGroup prepares the DeploymentsListByResourceGroup request.
func (c NginxDeploymentClient) preparerForDeploymentsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Nginx.NginxPlus/nginxDeployments", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForDeploymentsListByResourceGroupWithNextLink prepares the DeploymentsListByResourceGroup request with the given nextLink token.
func (c NginxDeploymentClient) preparerForDeploymentsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
	uri, err := url.Parse(nextLink)
	if err != nil {
		return nil, fmt.Errorf("parsing nextLink %q: %+v", nextLink, err)
	}
	queryParameters := map[string]interface{}{}
	for k, v := range uri.Query() {
		if len(v) == 0 {
			continue
		}
		val := v[0]
		val = autorest.Encode("query", val)
		queryParameters[k] = val
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeploymentsListByResourceGroup handles the response to the DeploymentsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c NginxDeploymentClient) responderForDeploymentsListByResourceGroup(resp *http.Response) (result DeploymentsListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []NginxDeployment `json:"value"`
		NextLink *string           `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	result.Model = &respObj.Values
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DeploymentsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForDeploymentsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDeploymentsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// DeploymentsListByResourceGroupComplete retrieves all of the results into a single object
func (c NginxDeploymentClient) DeploymentsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (DeploymentsListByResourceGroupCompleteResult, error) {
	return c.DeploymentsListByResourceGroupCompleteMatchingPredicate(ctx, id, NginxDeploymentOperationPredicate{})
}

// DeploymentsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c NginxDeploymentClient) DeploymentsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate NginxDeploymentOperationPredicate) (resp DeploymentsListByResourceGroupCompleteResult, err error) {
	items := make([]NginxDeployment, 0)

	page, err := c.DeploymentsListByResourceGroup(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading the initial page: %+v", err)
		return
	}
	if page.Model != nil {
		for _, v := range *page.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	for page.HasMore() {
		page, err = page.LoadMore(ctx)
		if err != nil {
			err = fmt.Errorf("loading the next page: %+v", err)
			return
		}

		if page.Model != nil {
			for _, v := range *page.Model {
				if predicate.Matches(v) {
					items = append(items, v)
				}
			}
		}
	}

	out := DeploymentsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
