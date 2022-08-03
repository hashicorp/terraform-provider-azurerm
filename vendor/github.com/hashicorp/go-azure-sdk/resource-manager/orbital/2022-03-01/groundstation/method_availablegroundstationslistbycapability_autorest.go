package groundstation

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

type AvailableGroundStationsListByCapabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AvailableGroundStation

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AvailableGroundStationsListByCapabilityOperationResponse, error)
}

type AvailableGroundStationsListByCapabilityCompleteResult struct {
	Items []AvailableGroundStation
}

func (r AvailableGroundStationsListByCapabilityOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AvailableGroundStationsListByCapabilityOperationResponse) LoadMore(ctx context.Context) (resp AvailableGroundStationsListByCapabilityOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type AvailableGroundStationsListByCapabilityOperationOptions struct {
	Capability *CapabilityParameter
}

func DefaultAvailableGroundStationsListByCapabilityOperationOptions() AvailableGroundStationsListByCapabilityOperationOptions {
	return AvailableGroundStationsListByCapabilityOperationOptions{}
}

func (o AvailableGroundStationsListByCapabilityOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o AvailableGroundStationsListByCapabilityOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Capability != nil {
		out["capability"] = *o.Capability
	}

	return out
}

// AvailableGroundStationsListByCapability ...
func (c GroundStationClient) AvailableGroundStationsListByCapability(ctx context.Context, id commonids.SubscriptionId, options AvailableGroundStationsListByCapabilityOperationOptions) (resp AvailableGroundStationsListByCapabilityOperationResponse, err error) {
	req, err := c.preparerForAvailableGroundStationsListByCapability(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "groundstation.GroundStationClient", "AvailableGroundStationsListByCapability", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "groundstation.GroundStationClient", "AvailableGroundStationsListByCapability", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAvailableGroundStationsListByCapability(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "groundstation.GroundStationClient", "AvailableGroundStationsListByCapability", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// AvailableGroundStationsListByCapabilityComplete retrieves all of the results into a single object
func (c GroundStationClient) AvailableGroundStationsListByCapabilityComplete(ctx context.Context, id commonids.SubscriptionId, options AvailableGroundStationsListByCapabilityOperationOptions) (AvailableGroundStationsListByCapabilityCompleteResult, error) {
	return c.AvailableGroundStationsListByCapabilityCompleteMatchingPredicate(ctx, id, options, AvailableGroundStationOperationPredicate{})
}

// AvailableGroundStationsListByCapabilityCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c GroundStationClient) AvailableGroundStationsListByCapabilityCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options AvailableGroundStationsListByCapabilityOperationOptions, predicate AvailableGroundStationOperationPredicate) (resp AvailableGroundStationsListByCapabilityCompleteResult, err error) {
	items := make([]AvailableGroundStation, 0)

	page, err := c.AvailableGroundStationsListByCapability(ctx, id, options)
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

	out := AvailableGroundStationsListByCapabilityCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForAvailableGroundStationsListByCapability prepares the AvailableGroundStationsListByCapability request.
func (c GroundStationClient) preparerForAvailableGroundStationsListByCapability(ctx context.Context, id commonids.SubscriptionId, options AvailableGroundStationsListByCapabilityOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Orbital/availableGroundStations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAvailableGroundStationsListByCapabilityWithNextLink prepares the AvailableGroundStationsListByCapability request with the given nextLink token.
func (c GroundStationClient) preparerForAvailableGroundStationsListByCapabilityWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAvailableGroundStationsListByCapability handles the response to the AvailableGroundStationsListByCapability request. The method always
// closes the http.Response Body.
func (c GroundStationClient) responderForAvailableGroundStationsListByCapability(resp *http.Response) (result AvailableGroundStationsListByCapabilityOperationResponse, err error) {
	type page struct {
		Values   []AvailableGroundStation `json:"value"`
		NextLink *string                  `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AvailableGroundStationsListByCapabilityOperationResponse, err error) {
			req, err := c.preparerForAvailableGroundStationsListByCapabilityWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "groundstation.GroundStationClient", "AvailableGroundStationsListByCapability", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "groundstation.GroundStationClient", "AvailableGroundStationsListByCapability", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAvailableGroundStationsListByCapability(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "groundstation.GroundStationClient", "AvailableGroundStationsListByCapability", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
