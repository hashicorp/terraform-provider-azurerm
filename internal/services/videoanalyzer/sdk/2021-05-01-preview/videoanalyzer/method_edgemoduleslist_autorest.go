package videoanalyzer

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type EdgeModulesListResponse struct {
	HttpResponse *http.Response
	Model        *[]EdgeModuleEntity

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (EdgeModulesListResponse, error)
}

type EdgeModulesListCompleteResult struct {
	Items []EdgeModuleEntity
}

func (r EdgeModulesListResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r EdgeModulesListResponse) LoadMore(ctx context.Context) (resp EdgeModulesListResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type EdgeModulesListOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultEdgeModulesListOptions() EdgeModulesListOptions {
	return EdgeModulesListOptions{}
}

func (o EdgeModulesListOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// EdgeModulesList ...
func (c VideoAnalyzerClient) EdgeModulesList(ctx context.Context, id VideoAnalyzerId, options EdgeModulesListOptions) (resp EdgeModulesListResponse, err error) {
	req, err := c.preparerForEdgeModulesList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForEdgeModulesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// EdgeModulesListComplete retrieves all of the results into a single object
func (c VideoAnalyzerClient) EdgeModulesListComplete(ctx context.Context, id VideoAnalyzerId, options EdgeModulesListOptions) (EdgeModulesListCompleteResult, error) {
	return c.EdgeModulesListCompleteMatchingPredicate(ctx, id, options, EdgeModuleEntityPredicate{})
}

// EdgeModulesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VideoAnalyzerClient) EdgeModulesListCompleteMatchingPredicate(ctx context.Context, id VideoAnalyzerId, options EdgeModulesListOptions, predicate EdgeModuleEntityPredicate) (resp EdgeModulesListCompleteResult, err error) {
	items := make([]EdgeModuleEntity, 0)

	page, err := c.EdgeModulesList(ctx, id, options)
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

	out := EdgeModulesListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForEdgeModulesList prepares the EdgeModulesList request.
func (c VideoAnalyzerClient) preparerForEdgeModulesList(ctx context.Context, id VideoAnalyzerId, options EdgeModulesListOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/edgeModules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForEdgeModulesListWithNextLink prepares the EdgeModulesList request with the given nextLink token.
func (c VideoAnalyzerClient) preparerForEdgeModulesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForEdgeModulesList handles the response to the EdgeModulesList request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForEdgeModulesList(resp *http.Response) (result EdgeModulesListResponse, err error) {
	type page struct {
		Values   []EdgeModuleEntity `json:"value"`
		NextLink *string            `json:"@nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result EdgeModulesListResponse, err error) {
			req, err := c.preparerForEdgeModulesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForEdgeModulesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
