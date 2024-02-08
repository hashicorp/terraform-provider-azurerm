package managedenvironments

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

type ManagedCertificatesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ManagedCertificate

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ManagedCertificatesListOperationResponse, error)
}

type ManagedCertificatesListCompleteResult struct {
	Items []ManagedCertificate
}

func (r ManagedCertificatesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ManagedCertificatesListOperationResponse) LoadMore(ctx context.Context) (resp ManagedCertificatesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ManagedCertificatesList ...
func (c ManagedEnvironmentsClient) ManagedCertificatesList(ctx context.Context, id ManagedEnvironmentId) (resp ManagedCertificatesListOperationResponse, err error) {
	req, err := c.preparerForManagedCertificatesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForManagedCertificatesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForManagedCertificatesList prepares the ManagedCertificatesList request.
func (c ManagedEnvironmentsClient) preparerForManagedCertificatesList(ctx context.Context, id ManagedEnvironmentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/managedCertificates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForManagedCertificatesListWithNextLink prepares the ManagedCertificatesList request with the given nextLink token.
func (c ManagedEnvironmentsClient) preparerForManagedCertificatesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForManagedCertificatesList handles the response to the ManagedCertificatesList request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForManagedCertificatesList(resp *http.Response) (result ManagedCertificatesListOperationResponse, err error) {
	type page struct {
		Values   []ManagedCertificate `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ManagedCertificatesListOperationResponse, err error) {
			req, err := c.preparerForManagedCertificatesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForManagedCertificatesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ManagedCertificatesListComplete retrieves all of the results into a single object
func (c ManagedEnvironmentsClient) ManagedCertificatesListComplete(ctx context.Context, id ManagedEnvironmentId) (ManagedCertificatesListCompleteResult, error) {
	return c.ManagedCertificatesListCompleteMatchingPredicate(ctx, id, ManagedCertificateOperationPredicate{})
}

// ManagedCertificatesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedEnvironmentsClient) ManagedCertificatesListCompleteMatchingPredicate(ctx context.Context, id ManagedEnvironmentId, predicate ManagedCertificateOperationPredicate) (resp ManagedCertificatesListCompleteResult, err error) {
	items := make([]ManagedCertificate, 0)

	page, err := c.ManagedCertificatesList(ctx, id)
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

	out := ManagedCertificatesListCompleteResult{
		Items: items,
	}
	return out, nil
}
