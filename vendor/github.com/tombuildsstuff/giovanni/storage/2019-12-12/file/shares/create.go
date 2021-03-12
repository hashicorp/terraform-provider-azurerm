package shares

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

type CreateInput struct {
	// Specifies the maximum size of the share, in gigabytes.
	// Must be greater than 0, and less than or equal to 5TB (5120).
	QuotaInGB int

	MetaData map[string]string
}

// Create creates the specified Storage Share within the specified Storage Account
func (client Client) Create(ctx context.Context, accountName, shareName string, input CreateInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "Create", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "Create", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "Create", "`shareName` must be a lower-cased string.")
	}
	if input.QuotaInGB <= 0 || input.QuotaInGB > 102400 {
		return result, validation.NewError("shares.Client", "Create", "`input.QuotaInGB` must be greater than 0, and less than/equal to 100TB (102400 GB)")
	}
	if err := metadata.Validate(input.MetaData); err != nil {
		return result, validation.NewError("shares.Client", "Create", fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
	}

	req, err := client.CreatePreparer(ctx, accountName, shareName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "Create", resp, "Failure responding to request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client Client) CreatePreparer(ctx context.Context, accountName, shareName string, input CreateInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype": autorest.Encode("path", "share"),
	}

	headers := map[string]interface{}{
		"x-ms-version":     APIVersion,
		"x-ms-share-quota": input.QuotaInGB,
	}

	headers = metadata.SetIntoHeaders(headers, input.MetaData)

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client Client) CreateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client Client) CreateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
