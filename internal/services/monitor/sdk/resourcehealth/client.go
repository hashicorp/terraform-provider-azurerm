package resourcehealth

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

const (
	// DefaultBaseURI is the default URI used for the service Resourcehealth
	DefaultBaseURI = "https://management.azure.com"
)

// BaseClient is the base client for Resourcehealth.
type BaseClient struct {
	autorest.Client
	BaseURI        string
	SubscriptionID string
}

type ResourceHealthMetadataClient struct {
	BaseClient
}

func NewResourceHealthMetadataClientWithBaseURI(baseURI string, subscriptionID string) ResourceHealthMetadataClient {
	return ResourceHealthMetadataClient{BaseClient: BaseClient{
		Client:         autorest.NewClientWithUserAgent(""),
		BaseURI:        baseURI,
		SubscriptionID: subscriptionID,
	}}
}

func (client ResourceHealthMetadataClient) GetMetaData(ctx context.Context) (result ResourceHealthMetadataResource, err error) {
	req, err := client.GetPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers/Microsoft.ResourceHealthMetadataClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "providers/Microsoft.ResourceHealthMetadataClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers/Microsoft.ResourceHealthMetadataClient", "Get", resp, "Failure responding to request")
		return
	}

	return result, nil
}

// GetPreparer prepares the Get request.
func (client ResourceHealthMetadataClient) GetPreparer(ctx context.Context) (*http.Request, error) {
	const APIVersion = "2018-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/providers/Microsoft.ResourceHealth/metadata"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (client ResourceHealthMetadataClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

func (client ResourceHealthMetadataClient) GetResponder(resp *http.Response) (result ResourceHealthMetadataResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
