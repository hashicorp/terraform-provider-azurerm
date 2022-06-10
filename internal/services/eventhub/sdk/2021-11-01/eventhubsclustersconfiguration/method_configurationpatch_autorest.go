package eventhubsclustersconfiguration

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConfigurationPatchOperationResponse struct {
	HttpResponse *http.Response
	Model        *ClusterQuotaConfigurationProperties
}

// ConfigurationPatch ...
func (c EventHubsClustersConfigurationClient) ConfigurationPatch(ctx context.Context, id ClusterId, input ClusterQuotaConfigurationProperties) (result ConfigurationPatchOperationResponse, err error) {
	req, err := c.preparerForConfigurationPatch(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersconfiguration.EventHubsClustersConfigurationClient", "ConfigurationPatch", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersconfiguration.EventHubsClustersConfigurationClient", "ConfigurationPatch", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConfigurationPatch(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersconfiguration.EventHubsClustersConfigurationClient", "ConfigurationPatch", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConfigurationPatch prepares the ConfigurationPatch request.
func (c EventHubsClustersConfigurationClient) preparerForConfigurationPatch(ctx context.Context, id ClusterId, input ClusterQuotaConfigurationProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/quotaConfiguration/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConfigurationPatch handles the response to the ConfigurationPatch request. The method always
// closes the http.Response Body.
func (c EventHubsClustersConfigurationClient) responderForConfigurationPatch(resp *http.Response) (result ConfigurationPatchOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
