package network

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

func getProfilePackageResult(client network.VirtualNetworkGatewaysClient, future network.VirtualNetworkGatewaysGetVpnProfilePackageURLFuture) (str string, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewaysGetVpnProfilePackageURLFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("network.VirtualNetworkGatewaysGetVpnProfilePackageURLFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	var s network.String
	if s.Response.Response, err = future.GetResult(sender); err == nil && s.Response.Response.StatusCode != http.StatusNoContent {
		str, err = getvpnprofilepackageResponder(client, s.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewaysGetVpnProfilePackageURLFuture", "Result", s.Response.Response, "Failure responding to request")
		}
	}
	return
}

func getvpnprofilepackageResponder(client network.VirtualNetworkGatewaysClient, resp *http.Response) (result string, err error) {
	byteArr := make([]byte, 1024)
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingBytes(&byteArr),
		autorest.ByClosing())
	result = string(byteArr[1 : len(byteArr)-1])
	return
}
