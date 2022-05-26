package clusterextensions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ExtensionsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Extension

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ExtensionsListOperationResponse, error)
}

type ExtensionsListCompleteResult struct {
	Items []Extension
}

func (r ExtensionsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ExtensionsListOperationResponse) LoadMore(ctx context.Context) (resp ExtensionsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ExtensionsList ...
func (c ClusterExtensionsClient) ExtensionsList(ctx context.Context, id ProviderId) (resp ExtensionsListOperationResponse, err error) {
	req, err := c.preparerForExtensionsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForExtensionsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ExtensionsListComplete retrieves all of the results into a single object
func (c ClusterExtensionsClient) ExtensionsListComplete(ctx context.Context, id ProviderId) (ExtensionsListCompleteResult, error) {
	return c.ExtensionsListCompleteMatchingPredicate(ctx, id, ExtensionOperationPredicate{})
}

// ExtensionsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ClusterExtensionsClient) ExtensionsListCompleteMatchingPredicate(ctx context.Context, id ProviderId, predicate ExtensionOperationPredicate) (resp ExtensionsListCompleteResult, err error) {
	items := make([]Extension, 0)

	page, err := c.ExtensionsList(ctx, id)
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

	out := ExtensionsListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForExtensionsList prepares the ExtensionsList request.
func (c ClusterExtensionsClient) preparerForExtensionsList(ctx context.Context, id ProviderId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.KubernetesConfiguration/extensions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForExtensionsListWithNextLink prepares the ExtensionsList request with the given nextLink token.
func (c ClusterExtensionsClient) preparerForExtensionsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForExtensionsList handles the response to the ExtensionsList request. The method always
// closes the http.Response Body.
func (c ClusterExtensionsClient) responderForExtensionsList(resp *http.Response) (result ExtensionsListOperationResponse, err error) {
	type page struct {
		Values   []Extension `json:"value"`
		NextLink *string     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ExtensionsListOperationResponse, err error) {
			req, err := c.preparerForExtensionsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForExtensionsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
