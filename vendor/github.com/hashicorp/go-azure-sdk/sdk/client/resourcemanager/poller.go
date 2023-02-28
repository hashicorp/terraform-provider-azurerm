package resourcemanager

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

func PollerFromResponse(response *client.Response, client *Client) (poller pollers.Poller, err error) {
	if response == nil {
		return pollers.Poller{}, fmt.Errorf("no HTTP Response was returned")
	}

	// If this is a LRO we should either have a 201/202 with a Polling URI header
	isLroStatus := response.StatusCode == http.StatusCreated || response.StatusCode == http.StatusAccepted
	lroPollingUri := pollingUriForLongRunningOperation(response)
	if isLroStatus && lroPollingUri != "" {
		lro, lroErr := longRunningOperationPollerFromResponse(response, client.Client)
		if lroErr != nil {
			err = lroErr
			return pollers.Poller{}, fmt.Errorf("building long-running-operation poller: %+v", lroErr)
		}
		return pollers.NewPoller(lro, lro.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
	}

	// or we should be polling on the `provisioningState` of the resource
	contentType := response.Header.Get("Content-Type")
	if contentType == "" && response.Request != nil {
		contentType = response.Request.Header.Get("Content-Type")
	}

	statusCodesToCheckProvisioningState := response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated
	contentTypeMatchesForProvisioningStateCheck := strings.Contains(strings.ToLower(contentType), "application/json")
	methodIsApplicable := strings.EqualFold(response.Request.Method, "PATCH") ||
		strings.EqualFold(response.Request.Method, "POST") ||
		strings.EqualFold(response.Request.Method, "PUT")
	if statusCodesToCheckProvisioningState && contentTypeMatchesForProvisioningStateCheck && methodIsApplicable {
		provisioningState, provisioningStateErr := provisioningStatePollerFromResponse(response, client, DefaultPollingInterval)
		if provisioningStateErr != nil {
			err = provisioningStateErr
			return pollers.Poller{}, fmt.Errorf("building provisioningState poller: %+v", provisioningStateErr)
		}
		return pollers.NewPoller(provisioningState, provisioningState.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
	}

	// finally, if it was a Delete that returned a 200/204
	methodIsDelete := strings.EqualFold(response.Request.Method, "DELETE")
	statusCodesToCheckDelete := response.StatusCode == http.StatusOK || response.StatusCode == http.StatusNoContent
	if methodIsDelete && statusCodesToCheckDelete {
		deletePoller, deletePollerErr := deletePollerFromResponse(response, client, DefaultPollingInterval)
		if deletePollerErr != nil {
			err = deletePollerErr
			return pollers.Poller{}, fmt.Errorf("building delete poller: %+v", deletePollerErr)
		}
		return pollers.NewPoller(deletePoller, deletePoller.initialRetryDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow), nil
	}

	return pollers.Poller{}, fmt.Errorf("no applicable pollers were found for the response")
}
