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

type CertificatesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Certificate

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (CertificatesListOperationResponse, error)
}

type CertificatesListCompleteResult struct {
	Items []Certificate
}

func (r CertificatesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r CertificatesListOperationResponse) LoadMore(ctx context.Context) (resp CertificatesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// CertificatesList ...
func (c ManagedEnvironmentsClient) CertificatesList(ctx context.Context, id ManagedEnvironmentId) (resp CertificatesListOperationResponse, err error) {
	req, err := c.preparerForCertificatesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForCertificatesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForCertificatesList prepares the CertificatesList request.
func (c ManagedEnvironmentsClient) preparerForCertificatesList(ctx context.Context, id ManagedEnvironmentId) (*http.Request, error) {
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

// preparerForCertificatesListWithNextLink prepares the CertificatesList request with the given nextLink token.
func (c ManagedEnvironmentsClient) preparerForCertificatesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForCertificatesList handles the response to the CertificatesList request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForCertificatesList(resp *http.Response) (result CertificatesListOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result CertificatesListOperationResponse, err error) {
			req, err := c.preparerForCertificatesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForCertificatesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// CertificatesListComplete retrieves all of the results into a single object
func (c ManagedEnvironmentsClient) CertificatesListComplete(ctx context.Context, id ManagedEnvironmentId) (CertificatesListCompleteResult, error) {
	return c.CertificatesListCompleteMatchingPredicate(ctx, id, CertificateOperationPredicate{})
}

// CertificatesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedEnvironmentsClient) CertificatesListCompleteMatchingPredicate(ctx context.Context, id ManagedEnvironmentId, predicate CertificateOperationPredicate) (resp CertificatesListCompleteResult, err error) {
	items := make([]Certificate, 0)

	page, err := c.CertificatesList(ctx, id)
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

	out := CertificatesListCompleteResult{
		Items: items,
	}
	return out, nil
}
