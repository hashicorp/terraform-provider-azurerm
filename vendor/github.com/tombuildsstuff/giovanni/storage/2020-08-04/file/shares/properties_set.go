package shares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// SetProperties lets you update the Quota for the specified Storage Share
func (client Client) SetProperties(ctx context.Context, accountName, shareName string, newQuotaGB int) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "SetProperties", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "SetProperties", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "SetProperties", "`shareName` must be a lower-cased string.")
	}
	if newQuotaGB <= 0 || newQuotaGB > 102400 {
		return result, validation.NewError("shares.Client", "SetProperties", "`newQuotaGB` must be greater than 0, and less than/equal to 100TB (102400 GB)")
	}

	req, err := client.SetPropertiesPreparer(ctx, accountName, shareName, newQuotaGB)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "SetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetPropertiesSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "SetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "SetProperties", resp, "Failure responding to request")
		return
	}

	return
}

// SetPropertiesPreparer prepares the SetProperties request.
func (client Client) SetPropertiesPreparer(ctx context.Context, accountName, shareName string, quotaGB int) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype": autorest.Encode("query", "share"),
		"comp":    autorest.Encode("query", "properties"),
	}

	headers := map[string]interface{}{
		"x-ms-version":     APIVersion,
		"x-ms-share-quota": quotaGB,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetPropertiesSender sends the SetProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetPropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetPropertiesResponder handles the response to the SetProperties request. The method always
// closes the http.Response Body.
func (client Client) SetPropertiesResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
