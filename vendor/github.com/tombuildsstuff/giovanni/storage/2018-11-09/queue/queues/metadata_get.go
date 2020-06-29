package queues

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetMetaDataResult struct {
	autorest.Response

	MetaData map[string]string
}

// GetMetaData returns the metadata for this Queue
func (client Client) GetMetaData(ctx context.Context, accountName, queueName string) (result GetMetaDataResult, err error) {
	if accountName == "" {
		return result, validation.NewError("queues.Client", "GetMetaData", "`accountName` cannot be an empty string.")
	}
	if queueName == "" {
		return result, validation.NewError("queues.Client", "GetMetaData", "`queueName` cannot be an empty string.")
	}
	if strings.ToLower(queueName) != queueName {
		return result, validation.NewError("queues.Client", "GetMetaData", "`queueName` must be a lower-cased string.")
	}

	req, err := client.GetMetaDataPreparer(ctx, accountName, queueName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queues.Client", "GetMetaData", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetMetaDataSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "queues.Client", "GetMetaData", resp, "Failure sending request")
		return
	}

	result, err = client.GetMetaDataResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queues.Client", "GetMetaData", resp, "Failure responding to request")
		return
	}

	return
}

// GetMetaDataPreparer prepares the GetMetaData request.
func (client Client) GetMetaDataPreparer(ctx context.Context, accountName, queueName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"queueName": autorest.Encode("path", queueName),
	}

	queryParameters := map[string]interface{}{
		"comp": autorest.Encode("path", "metadata"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetQueueEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{queueName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetMetaDataSender sends the GetMetaData request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetMetaDataSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetMetaDataResponder handles the response to the GetMetaData request. The method always
// closes the http.Response Body.
func (client Client) GetMetaDataResponder(resp *http.Response) (result GetMetaDataResult, err error) {
	if resp != nil {
		result.MetaData = metadata.ParseFromHeaders(resp.Header)
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
