package azuresdkhacks

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/Azure/go-autorest/autorest"
)

// UpdateNetworkInterfaceAllowingRemovalOfNSG patches our way around a design flaw in the Azure
// Resource Manager API <-> Azure SDK for Go where it's not possible to remove a Network Security Group
func UpdateNetworkInterfaceAllowingRemovalOfNSG(ctx context.Context, client *network.InterfacesClient, resourceGroupName string, networkInterfaceName string, parameters network.Interface) (result network.InterfacesCreateOrUpdateFuture, err error) {
	req, err := updateNetworkInterfaceAllowingRemovalOfNSGPreparer(ctx, client, resourceGroupName, networkInterfaceName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.InterfacesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.InterfacesClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// updateNetworkInterfaceAllowingRemovalOfNSGPreparer prepares the CreateOrUpdate request but applies the
// necessary patches to be able to remove the NSG if required
func updateNetworkInterfaceAllowingRemovalOfNSGPreparer(ctx context.Context, client *network.InterfacesClient, resourceGroupName string, networkInterfaceName string, parameters network.Interface) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"networkInterfaceName": autorest.Encode("path", networkInterfaceName),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	parameters.Etag = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkInterfaces/{networkInterfaceName}", pathParameters),
		withJsonWorkingAroundTheBrokenNetworkAPI(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func withJsonWorkingAroundTheBrokenNetworkAPI(v network.Interface) autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err == nil {
				b, err := json.Marshal(v)
				if err == nil {
					// there's a few fields which can be intentionally set to nil - as such here we need to check if they should be nil and force them to be nil
					var out map[string]interface{}
					if err := json.Unmarshal(b, &out); err != nil {
						return r, err
					}

					// apply the hack
					out = patchNICUpdateAPIIssue(v, out)

					// then reserialize it as needed
					b, err = json.Marshal(out)
					if err == nil {
						r.ContentLength = int64(len(b))
						r.Body = ioutil.NopCloser(bytes.NewReader(b))
					}
				}
			}
			return r, err
		})
	}
}

func patchNICUpdateAPIIssue(nic network.Interface, input map[string]interface{}) map[string]interface{} {
	if nic.InterfacePropertiesFormat == nil {
		return input
	}

	output := input

	if v, ok := output["properties"]; ok {
		props := v.(map[string]interface{})

		if nic.InterfacePropertiesFormat.NetworkSecurityGroup == nil {
			var hack *string // a nil-pointered string
			props["networkSecurityGroup"] = hack
		}

		output["properties"] = props
	}

	return output
}
