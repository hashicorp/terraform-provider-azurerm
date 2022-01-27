package profiles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GenerateSsoUriResponse struct {
	HttpResponse *http.Response
	Model        *SsoUri
}

// GenerateSsoUri ...
func (c ProfilesClient) GenerateSsoUri(ctx context.Context, id ProfileId) (result GenerateSsoUriResponse, err error) {
	req, err := c.preparerForGenerateSsoUri(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "GenerateSsoUri", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "GenerateSsoUri", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGenerateSsoUri(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "GenerateSsoUri", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGenerateSsoUri prepares the GenerateSsoUri request.
func (c ProfilesClient) preparerForGenerateSsoUri(ctx context.Context, id ProfileId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/generateSsoUri", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGenerateSsoUri handles the response to the GenerateSsoUri request. The method always
// closes the http.Response Body.
func (c ProfilesClient) responderForGenerateSsoUri(resp *http.Response) (result GenerateSsoUriResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
