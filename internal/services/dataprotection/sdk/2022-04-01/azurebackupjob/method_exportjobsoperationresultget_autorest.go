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

type ExportJobsOperationResultGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ExportJobsResult
}

// ExportJobsOperationResultGet ...
func (c AzureBackupJobClient) ExportJobsOperationResultGet(ctx context.Context, id OperationIdId) (result ExportJobsOperationResultGetOperationResponse, err error) {
	req, err := c.preparerForExportJobsOperationResultGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "azurebackupjob.AzureBackupJobClient", "ExportJobsOperationResultGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "azurebackupjob.AzureBackupJobClient", "ExportJobsOperationResultGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForExportJobsOperationResultGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "azurebackupjob.AzureBackupJobClient", "ExportJobsOperationResultGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForExportJobsOperationResultGet prepares the ExportJobsOperationResultGet request.
func (c AzureBackupJobClient) preparerForExportJobsOperationResultGet(ctx context.Context, id OperationIdId) (*http.Request, error) {
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

// responderForExportJobsOperationResultGet handles the response to the ExportJobsOperationResultGet request. The method always
// closes the http.Response Body.
func (c AzureBackupJobClient) responderForExportJobsOperationResultGet(resp *http.Response) (result ExportJobsOperationResultGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
