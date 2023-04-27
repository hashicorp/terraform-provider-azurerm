package iotconnectors

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

type FhirDestinationsListByIotConnectorOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]IotFhirDestination

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (FhirDestinationsListByIotConnectorOperationResponse, error)
}

type FhirDestinationsListByIotConnectorCompleteResult struct {
	Items []IotFhirDestination
}

func (r FhirDestinationsListByIotConnectorOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r FhirDestinationsListByIotConnectorOperationResponse) LoadMore(ctx context.Context) (resp FhirDestinationsListByIotConnectorOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// FhirDestinationsListByIotConnector ...
func (c IotConnectorsClient) FhirDestinationsListByIotConnector(ctx context.Context, id IotConnectorId) (resp FhirDestinationsListByIotConnectorOperationResponse, err error) {
	req, err := c.preparerForFhirDestinationsListByIotConnector(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "FhirDestinationsListByIotConnector", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "FhirDestinationsListByIotConnector", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForFhirDestinationsListByIotConnector(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "FhirDestinationsListByIotConnector", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForFhirDestinationsListByIotConnector prepares the FhirDestinationsListByIotConnector request.
func (c IotConnectorsClient) preparerForFhirDestinationsListByIotConnector(ctx context.Context, id IotConnectorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/fhirDestinations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForFhirDestinationsListByIotConnectorWithNextLink prepares the FhirDestinationsListByIotConnector request with the given nextLink token.
func (c IotConnectorsClient) preparerForFhirDestinationsListByIotConnectorWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForFhirDestinationsListByIotConnector handles the response to the FhirDestinationsListByIotConnector request. The method always
// closes the http.Response Body.
func (c IotConnectorsClient) responderForFhirDestinationsListByIotConnector(resp *http.Response) (result FhirDestinationsListByIotConnectorOperationResponse, err error) {
	type page struct {
		Values   []IotFhirDestination `json:"value"`
		NextLink *string              `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result FhirDestinationsListByIotConnectorOperationResponse, err error) {
			req, err := c.preparerForFhirDestinationsListByIotConnectorWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "FhirDestinationsListByIotConnector", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "FhirDestinationsListByIotConnector", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForFhirDestinationsListByIotConnector(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "FhirDestinationsListByIotConnector", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// FhirDestinationsListByIotConnectorComplete retrieves all of the results into a single object
func (c IotConnectorsClient) FhirDestinationsListByIotConnectorComplete(ctx context.Context, id IotConnectorId) (FhirDestinationsListByIotConnectorCompleteResult, error) {
	return c.FhirDestinationsListByIotConnectorCompleteMatchingPredicate(ctx, id, IotFhirDestinationOperationPredicate{})
}

// FhirDestinationsListByIotConnectorCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c IotConnectorsClient) FhirDestinationsListByIotConnectorCompleteMatchingPredicate(ctx context.Context, id IotConnectorId, predicate IotFhirDestinationOperationPredicate) (resp FhirDestinationsListByIotConnectorCompleteResult, err error) {
	items := make([]IotFhirDestination, 0)

	page, err := c.FhirDestinationsListByIotConnector(ctx, id)
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

	out := FhirDestinationsListByIotConnectorCompleteResult{
		Items: items,
	}
	return out, nil
}
