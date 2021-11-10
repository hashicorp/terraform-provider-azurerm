package eventhubsclustersconfiguration

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConfigurationGetResponse struct {
	HttpResponse *http.Response
	Model        *ClusterQuotaConfigurationProperties
}

// ConfigurationGet ...
func (c EventHubsClustersConfigurationClient) ConfigurationGet(ctx context.Context, id ClusterId) (result ConfigurationGetResponse, err error) {
	req, err := c.preparerForConfigurationGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersconfiguration.EventHubsClustersConfigurationClient", "ConfigurationGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersconfiguration.EventHubsClustersConfigurationClient", "ConfigurationGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConfigurationGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersconfiguration.EventHubsClustersConfigurationClient", "ConfigurationGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConfigurationGet prepares the ConfigurationGet request.
func (c EventHubsClustersConfigurationClient) preparerForConfigurationGet(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/quotaConfiguration/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConfigurationGet handles the response to the ConfigurationGet request. The method always
// closes the http.Response Body.
func (c EventHubsClustersConfigurationClient) responderForConfigurationGet(resp *http.Response) (result ConfigurationGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
