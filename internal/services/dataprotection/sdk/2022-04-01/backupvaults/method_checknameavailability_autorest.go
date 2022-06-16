package backupvaults

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type CheckNameAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityResult
}

// CheckNameAvailability ...
func (c BackupVaultsClient) CheckNameAvailability(ctx context.Context, id ProviderLocationId, input CheckNameAvailabilityRequest) (result CheckNameAvailabilityOperationResponse, err error) {
	req, err := c.preparerForCheckNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "CheckNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "CheckNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "CheckNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckNameAvailability prepares the CheckNameAvailability request.
func (c BackupVaultsClient) preparerForCheckNameAvailability(ctx context.Context, id ProviderLocationId, input CheckNameAvailabilityRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckNameAvailability handles the response to the CheckNameAvailability request. The method always
// closes the http.Response Body.
func (c BackupVaultsClient) responderForCheckNameAvailability(resp *http.Response) (result CheckNameAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
