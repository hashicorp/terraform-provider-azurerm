package arcsettings

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcSettingsListByClusterOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ArcSetting

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ArcSettingsListByClusterOperationResponse, error)
}

type ArcSettingsListByClusterCompleteResult struct {
	Items []ArcSetting
}

func (r ArcSettingsListByClusterOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ArcSettingsListByClusterOperationResponse) LoadMore(ctx context.Context) (resp ArcSettingsListByClusterOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ArcSettingsListByCluster ...
func (c ArcSettingsClient) ArcSettingsListByCluster(ctx context.Context, id ClusterId) (resp ArcSettingsListByClusterOperationResponse, err error) {
	req, err := c.preparerForArcSettingsListByCluster(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsListByCluster", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsListByCluster", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForArcSettingsListByCluster(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsListByCluster", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForArcSettingsListByCluster prepares the ArcSettingsListByCluster request.
func (c ArcSettingsClient) preparerForArcSettingsListByCluster(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/arcSettings", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForArcSettingsListByClusterWithNextLink prepares the ArcSettingsListByCluster request with the given nextLink token.
func (c ArcSettingsClient) preparerForArcSettingsListByClusterWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForArcSettingsListByCluster handles the response to the ArcSettingsListByCluster request. The method always
// closes the http.Response Body.
func (c ArcSettingsClient) responderForArcSettingsListByCluster(resp *http.Response) (result ArcSettingsListByClusterOperationResponse, err error) {
	type page struct {
		Values   []ArcSetting `json:"value"`
		NextLink *string      `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ArcSettingsListByClusterOperationResponse, err error) {
			req, err := c.preparerForArcSettingsListByClusterWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsListByCluster", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsListByCluster", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForArcSettingsListByCluster(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsListByCluster", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ArcSettingsListByClusterComplete retrieves all of the results into a single object
func (c ArcSettingsClient) ArcSettingsListByClusterComplete(ctx context.Context, id ClusterId) (ArcSettingsListByClusterCompleteResult, error) {
	return c.ArcSettingsListByClusterCompleteMatchingPredicate(ctx, id, ArcSettingOperationPredicate{})
}

// ArcSettingsListByClusterCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ArcSettingsClient) ArcSettingsListByClusterCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate ArcSettingOperationPredicate) (resp ArcSettingsListByClusterCompleteResult, err error) {
	items := make([]ArcSetting, 0)

	page, err := c.ArcSettingsListByCluster(ctx, id)
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

	out := ArcSettingsListByClusterCompleteResult{
		Items: items,
	}
	return out, nil
}
