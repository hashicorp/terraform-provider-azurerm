package domaintopics

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByDomainResponse struct {
	HttpResponse *http.Response
	Model        *[]DomainTopic

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByDomainResponse, error)
}

type ListByDomainCompleteResult struct {
	Items []DomainTopic
}

func (r ListByDomainResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByDomainResponse) LoadMore(ctx context.Context) (resp ListByDomainResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByDomainOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListByDomainOptions() ListByDomainOptions {
	return ListByDomainOptions{}
}

func (o ListByDomainOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListByDomain ...
func (c DomainTopicsClient) ListByDomain(ctx context.Context, id DomainId, options ListByDomainOptions) (resp ListByDomainResponse, err error) {
	req, err := c.preparerForListByDomain(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "domaintopics.DomainTopicsClient", "ListByDomain", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "domaintopics.DomainTopicsClient", "ListByDomain", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByDomain(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "domaintopics.DomainTopicsClient", "ListByDomain", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByDomainComplete retrieves all of the results into a single object
func (c DomainTopicsClient) ListByDomainComplete(ctx context.Context, id DomainId, options ListByDomainOptions) (ListByDomainCompleteResult, error) {
	return c.ListByDomainCompleteMatchingPredicate(ctx, id, options, DomainTopicPredicate{})
}

// ListByDomainCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DomainTopicsClient) ListByDomainCompleteMatchingPredicate(ctx context.Context, id DomainId, options ListByDomainOptions, predicate DomainTopicPredicate) (resp ListByDomainCompleteResult, err error) {
	items := make([]DomainTopic, 0)

	page, err := c.ListByDomain(ctx, id, options)
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

	out := ListByDomainCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByDomain prepares the ListByDomain request.
func (c DomainTopicsClient) preparerForListByDomain(ctx context.Context, id DomainId, options ListByDomainOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/topics", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByDomainWithNextLink prepares the ListByDomain request with the given nextLink token.
func (c DomainTopicsClient) preparerForListByDomainWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByDomain handles the response to the ListByDomain request. The method always
// closes the http.Response Body.
func (c DomainTopicsClient) responderForListByDomain(resp *http.Response) (result ListByDomainResponse, err error) {
	type page struct {
		Values   []DomainTopic `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByDomainResponse, err error) {
			req, err := c.preparerForListByDomainWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "domaintopics.DomainTopicsClient", "ListByDomain", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "domaintopics.DomainTopicsClient", "ListByDomain", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByDomain(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "domaintopics.DomainTopicsClient", "ListByDomain", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
