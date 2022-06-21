package datasources

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByWorkspaceOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DataSource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByWorkspaceOperationResponse, error)
}

type ListByWorkspaceCompleteResult struct {
	Items []DataSource
}

func (r ListByWorkspaceOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByWorkspaceOperationResponse) LoadMore(ctx context.Context) (resp ListByWorkspaceOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByWorkspaceOperationOptions struct {
	Filter *string
}

func DefaultListByWorkspaceOperationOptions() ListByWorkspaceOperationOptions {
	return ListByWorkspaceOperationOptions{}
}

func (o ListByWorkspaceOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByWorkspaceOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListByWorkspace ...
func (c DataSourcesClient) ListByWorkspace(ctx context.Context, id WorkspaceId, options ListByWorkspaceOperationOptions) (resp ListByWorkspaceOperationResponse, err error) {
	req, err := c.preparerForListByWorkspace(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datasources.DataSourcesClient", "ListByWorkspace", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "datasources.DataSourcesClient", "ListByWorkspace", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByWorkspace(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datasources.DataSourcesClient", "ListByWorkspace", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByWorkspaceComplete retrieves all of the results into a single object
func (c DataSourcesClient) ListByWorkspaceComplete(ctx context.Context, id WorkspaceId, options ListByWorkspaceOperationOptions) (ListByWorkspaceCompleteResult, error) {
	return c.ListByWorkspaceCompleteMatchingPredicate(ctx, id, options, DataSourceOperationPredicate{})
}

// ListByWorkspaceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DataSourcesClient) ListByWorkspaceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options ListByWorkspaceOperationOptions, predicate DataSourceOperationPredicate) (resp ListByWorkspaceCompleteResult, err error) {
	items := make([]DataSource, 0)

	page, err := c.ListByWorkspace(ctx, id, options)
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

	out := ListByWorkspaceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByWorkspace prepares the ListByWorkspace request.
func (c DataSourcesClient) preparerForListByWorkspace(ctx context.Context, id WorkspaceId, options ListByWorkspaceOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/dataSources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByWorkspaceWithNextLink prepares the ListByWorkspace request with the given nextLink token.
func (c DataSourcesClient) preparerForListByWorkspaceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByWorkspace handles the response to the ListByWorkspace request. The method always
// closes the http.Response Body.
func (c DataSourcesClient) responderForListByWorkspace(resp *http.Response) (result ListByWorkspaceOperationResponse, err error) {
	type page struct {
		Values   []DataSource `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByWorkspaceOperationResponse, err error) {
			req, err := c.preparerForListByWorkspaceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "datasources.DataSourcesClient", "ListByWorkspace", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "datasources.DataSourcesClient", "ListByWorkspace", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByWorkspace(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "datasources.DataSourcesClient", "ListByWorkspace", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
