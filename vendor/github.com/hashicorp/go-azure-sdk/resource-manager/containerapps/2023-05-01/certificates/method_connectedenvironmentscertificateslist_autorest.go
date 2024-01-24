package certificates

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

type ConnectedEnvironmentsCertificatesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Certificate

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ConnectedEnvironmentsCertificatesListOperationResponse, error)
}

type ConnectedEnvironmentsCertificatesListCompleteResult struct {
	Items []Certificate
}

func (r ConnectedEnvironmentsCertificatesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ConnectedEnvironmentsCertificatesListOperationResponse) LoadMore(ctx context.Context) (resp ConnectedEnvironmentsCertificatesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ConnectedEnvironmentsCertificatesList ...
func (c CertificatesClient) ConnectedEnvironmentsCertificatesList(ctx context.Context, id ConnectedEnvironmentId) (resp ConnectedEnvironmentsCertificatesListOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsCertificatesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForConnectedEnvironmentsCertificatesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForConnectedEnvironmentsCertificatesList prepares the ConnectedEnvironmentsCertificatesList request.
func (c CertificatesClient) preparerForConnectedEnvironmentsCertificatesList(ctx context.Context, id ConnectedEnvironmentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/certificates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForConnectedEnvironmentsCertificatesListWithNextLink prepares the ConnectedEnvironmentsCertificatesList request with the given nextLink token.
func (c CertificatesClient) preparerForConnectedEnvironmentsCertificatesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForConnectedEnvironmentsCertificatesList handles the response to the ConnectedEnvironmentsCertificatesList request. The method always
// closes the http.Response Body.
func (c CertificatesClient) responderForConnectedEnvironmentsCertificatesList(resp *http.Response) (result ConnectedEnvironmentsCertificatesListOperationResponse, err error) {
	type page struct {
		Values   []Certificate `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ConnectedEnvironmentsCertificatesListOperationResponse, err error) {
			req, err := c.preparerForConnectedEnvironmentsCertificatesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForConnectedEnvironmentsCertificatesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ConnectedEnvironmentsCertificatesListComplete retrieves all of the results into a single object
func (c CertificatesClient) ConnectedEnvironmentsCertificatesListComplete(ctx context.Context, id ConnectedEnvironmentId) (ConnectedEnvironmentsCertificatesListCompleteResult, error) {
	return c.ConnectedEnvironmentsCertificatesListCompleteMatchingPredicate(ctx, id, CertificateOperationPredicate{})
}

// ConnectedEnvironmentsCertificatesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CertificatesClient) ConnectedEnvironmentsCertificatesListCompleteMatchingPredicate(ctx context.Context, id ConnectedEnvironmentId, predicate CertificateOperationPredicate) (resp ConnectedEnvironmentsCertificatesListCompleteResult, err error) {
	items := make([]Certificate, 0)

	page, err := c.ConnectedEnvironmentsCertificatesList(ctx, id)
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

	out := ConnectedEnvironmentsCertificatesListCompleteResult{
		Items: items,
	}
	return out, nil
}
