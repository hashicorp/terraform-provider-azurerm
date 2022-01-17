package roles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByServerGroupResponse struct {
	HttpResponse *http.Response
	Model        *RoleListResult
}

// ListByServerGroup ...
func (c RolesClient) ListByServerGroup(ctx context.Context, id ServerGroupsv2Id) (result ListByServerGroupResponse, err error) {
	req, err := c.preparerForListByServerGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "roles.RolesClient", "ListByServerGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "roles.RolesClient", "ListByServerGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListByServerGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "roles.RolesClient", "ListByServerGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListByServerGroup prepares the ListByServerGroup request.
func (c RolesClient) preparerForListByServerGroup(ctx context.Context, id ServerGroupsv2Id) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/roles", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByServerGroup handles the response to the ListByServerGroup request. The method always
// closes the http.Response Body.
func (c RolesClient) responderForListByServerGroup(resp *http.Response) (result ListByServerGroupResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
