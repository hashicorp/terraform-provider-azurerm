package collection

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type ServicesListByResourceGroupResponse struct {
	HttpResponse *http.Response
	Model        *[]ServicesDescription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ServicesListByResourceGroupResponse, error)
}

type ServicesListByResourceGroupCompleteResult struct {
	Items []ServicesDescription
}

func (r ServicesListByResourceGroupResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ServicesListByResourceGroupResponse) LoadMore(ctx context.Context) (resp ServicesListByResourceGroupResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ServicesListByResourceGroup ...
func (c CollectionClient) ServicesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp ServicesListByResourceGroupResponse, err error) {
	req, err := c.preparerForServicesListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForServicesListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ServicesListByResourceGroupComplete retrieves all of the results into a single object
func (c CollectionClient) ServicesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (ServicesListByResourceGroupCompleteResult, error) {
	return c.ServicesListByResourceGroupCompleteMatchingPredicate(ctx, id, ServicesDescriptionPredicate{})
}

// ServicesListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CollectionClient) ServicesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ServicesDescriptionPredicate) (resp ServicesListByResourceGroupCompleteResult, err error) {
	items := make([]ServicesDescription, 0)

	page, err := c.ServicesListByResourceGroup(ctx, id)
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

	out := ServicesListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForServicesListByResourceGroup prepares the ServicesListByResourceGroup request.
func (c CollectionClient) preparerForServicesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.HealthcareApis/services", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForServicesListByResourceGroupWithNextLink prepares the ServicesListByResourceGroup request with the given nextLink token.
func (c CollectionClient) preparerForServicesListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForServicesListByResourceGroup handles the response to the ServicesListByResourceGroup request. The method always
// closes the http.Response Body.
func (c CollectionClient) responderForServicesListByResourceGroup(resp *http.Response) (result ServicesListByResourceGroupResponse, err error) {
	type page struct {
		Values   []ServicesDescription `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ServicesListByResourceGroupResponse, err error) {
			req, err := c.preparerForServicesListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForServicesListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
