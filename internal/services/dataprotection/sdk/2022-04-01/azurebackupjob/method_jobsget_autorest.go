package azurebackupjob

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

type JobsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *AzureBackupJobResource
}

// JobsGet ...
func (c AzureBackupJobClient) JobsGet(ctx context.Context, id BackupJobId) (result JobsGetOperationResponse, err error) {
	req, err := c.preparerForJobsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "azurebackupjob.AzureBackupJobClient", "JobsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "azurebackupjob.AzureBackupJobClient", "JobsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForJobsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "azurebackupjob.AzureBackupJobClient", "JobsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForJobsGet prepares the JobsGet request.
func (c AzureBackupJobClient) preparerForJobsGet(ctx context.Context, id BackupJobId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForJobsGet handles the response to the JobsGet request. The method always
// closes the http.Response Body.
func (c AzureBackupJobClient) responderForJobsGet(resp *http.Response) (result JobsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
