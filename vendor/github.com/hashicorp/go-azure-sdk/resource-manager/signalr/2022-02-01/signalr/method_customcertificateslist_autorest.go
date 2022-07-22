package signalr

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

type CustomCertificatesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]CustomCertificate

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (CustomCertificatesListOperationResponse, error)
}

type CustomCertificatesListCompleteResult struct {
	Items []CustomCertificate
}

func (r CustomCertificatesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r CustomCertificatesListOperationResponse) LoadMore(ctx context.Context) (resp CustomCertificatesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// CustomCertificatesList ...
func (c SignalRClient) CustomCertificatesList(ctx context.Context, id SignalRId) (resp CustomCertificatesListOperationResponse, err error) {
	req, err := c.preparerForCustomCertificatesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForCustomCertificatesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// CustomCertificatesListComplete retrieves all of the results into a single object
func (c SignalRClient) CustomCertificatesListComplete(ctx context.Context, id SignalRId) (CustomCertificatesListCompleteResult, error) {
	return c.CustomCertificatesListCompleteMatchingPredicate(ctx, id, CustomCertificateOperationPredicate{})
}

// CustomCertificatesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SignalRClient) CustomCertificatesListCompleteMatchingPredicate(ctx context.Context, id SignalRId, predicate CustomCertificateOperationPredicate) (resp CustomCertificatesListCompleteResult, err error) {
	items := make([]CustomCertificate, 0)

	page, err := c.CustomCertificatesList(ctx, id)
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

	out := CustomCertificatesListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForCustomCertificatesList prepares the CustomCertificatesList request.
func (c SignalRClient) preparerForCustomCertificatesList(ctx context.Context, id SignalRId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/customCertificates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForCustomCertificatesListWithNextLink prepares the CustomCertificatesList request with the given nextLink token.
func (c SignalRClient) preparerForCustomCertificatesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForCustomCertificatesList handles the response to the CustomCertificatesList request. The method always
// closes the http.Response Body.
func (c SignalRClient) responderForCustomCertificatesList(resp *http.Response) (result CustomCertificatesListOperationResponse, err error) {
	type page struct {
		Values   []CustomCertificate `json:"value"`
		NextLink *string             `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result CustomCertificatesListOperationResponse, err error) {
			req, err := c.preparerForCustomCertificatesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForCustomCertificatesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
