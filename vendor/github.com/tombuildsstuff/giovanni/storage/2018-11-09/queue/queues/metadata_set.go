package queues

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

// SetMetaData returns the metadata for this Queue
func (client Client) SetMetaData(ctx context.Context, accountName, queueName string, metaData map[string]string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("queues.Client", "SetMetaData", "`accountName` cannot be an empty string.")
	}
	if queueName == "" {
		return result, validation.NewError("queues.Client", "SetMetaData", "`queueName` cannot be an empty string.")
	}
	if strings.ToLower(queueName) != queueName {
		return result, validation.NewError("queues.Client", "SetMetaData", "`queueName` must be a lower-cased string.")
	}
	if err := metadata.Validate(metaData); err != nil {
		return result, validation.NewError("queues.Client", "SetMetaData", fmt.Sprintf("`metadata` is not valid: %s.", err))
	}

	req, err := client.SetMetaDataPreparer(ctx, accountName, queueName, metaData)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queues.Client", "SetMetaData", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetMetaDataSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "queues.Client", "SetMetaData", resp, "Failure sending request")
		return
	}

	result, err = client.SetMetaDataResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queues.Client", "SetMetaData", resp, "Failure responding to request")
		return
	}

	return
}

// SetMetaDataPreparer prepares the SetMetaData request.
func (client Client) SetMetaDataPreparer(ctx context.Context, accountName, queueName string, metaData map[string]string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"queueName": autorest.Encode("path", queueName),
	}

	queryParameters := map[string]interface{}{
		"comp": autorest.Encode("path", "metadata"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	headers = metadata.SetIntoHeaders(headers, metaData)

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetQueueEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{queueName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetMetaDataSender sends the SetMetaData request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetMetaDataSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetMetaDataResponder handles the response to the SetMetaData request. The method always
// closes the http.Response Body.
func (client Client) SetMetaDataResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
