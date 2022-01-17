package configurations

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByServerGroupResponse struct {
	HttpResponse *http.Response
	Model        *[]ServerGroupConfiguration

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByServerGroupResponse, error)
}

type ListByServerGroupCompleteResult struct {
	Items []ServerGroupConfiguration
}

func (r ListByServerGroupResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByServerGroupResponse) LoadMore(ctx context.Context) (resp ListByServerGroupResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByServerGroup ...
func (c ConfigurationsClient) ListByServerGroup(ctx context.Context, id ServerGroupsv2Id) (resp ListByServerGroupResponse, err error) {
	req, err := c.preparerForListByServerGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServerGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServerGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByServerGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServerGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByServerGroupComplete retrieves all of the results into a single object
func (c ConfigurationsClient) ListByServerGroupComplete(ctx context.Context, id ServerGroupsv2Id) (ListByServerGroupCompleteResult, error) {
	return c.ListByServerGroupCompleteMatchingPredicate(ctx, id, ServerGroupConfigurationPredicate{})
}

// ListByServerGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ConfigurationsClient) ListByServerGroupCompleteMatchingPredicate(ctx context.Context, id ServerGroupsv2Id, predicate ServerGroupConfigurationPredicate) (resp ListByServerGroupCompleteResult, err error) {
	items := make([]ServerGroupConfiguration, 0)

	page, err := c.ListByServerGroup(ctx, id)
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

	out := ListByServerGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByServerGroup prepares the ListByServerGroup request.
func (c ConfigurationsClient) preparerForListByServerGroup(ctx context.Context, id ServerGroupsv2Id) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/configurations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByServerGroupWithNextLink prepares the ListByServerGroup request with the given nextLink token.
func (c ConfigurationsClient) preparerForListByServerGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByServerGroup handles the response to the ListByServerGroup request. The method always
// closes the http.Response Body.
func (c ConfigurationsClient) responderForListByServerGroup(resp *http.Response) (result ListByServerGroupResponse, err error) {
	type page struct {
		Values   []ServerGroupConfiguration `json:"value"`
		NextLink *string                    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByServerGroupResponse, err error) {
			req, err := c.preparerForListByServerGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServerGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServerGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByServerGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServerGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
