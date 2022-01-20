package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type DeletedAccountsPurgeResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DeletedAccountsPurge ...
func (c CognitiveServicesAccountsClient) DeletedAccountsPurge(ctx context.Context, id DeletedAccountId) (result DeletedAccountsPurgeResponse, err error) {
	req, err := c.preparerForDeletedAccountsPurge(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsPurge", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDeletedAccountsPurge(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsPurge", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeletedAccountsPurgeThenPoll performs DeletedAccountsPurge then polls until it's completed
func (c CognitiveServicesAccountsClient) DeletedAccountsPurgeThenPoll(ctx context.Context, id DeletedAccountId) error {
	result, err := c.DeletedAccountsPurge(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DeletedAccountsPurge: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DeletedAccountsPurge: %+v", err)
	}

	return nil
}

// preparerForDeletedAccountsPurge prepares the DeletedAccountsPurge request.
func (c CognitiveServicesAccountsClient) preparerForDeletedAccountsPurge(ctx context.Context, id DeletedAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDeletedAccountsPurge sends the DeletedAccountsPurge request. The method will close the
// http.Response Body if it receives an error.
func (c CognitiveServicesAccountsClient) senderForDeletedAccountsPurge(ctx context.Context, req *http.Request) (future DeletedAccountsPurgeResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
