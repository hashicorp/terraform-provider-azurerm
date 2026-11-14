// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
)

// functionAppFunctionDefaultKeyName is the name of the key Azure creates automatically for every Function.
const functionAppFunctionDefaultKeyName = "default"

// listFunctionAppFunctionKeys retrieves the Function-level keys (not the Host keys) for the specified Function.
func listFunctionAppFunctionKeys(ctx context.Context, client *webapps.WebAppsClient, id webapps.FunctionId) (map[string]string, error) {
	resp, err := client.ListFunctionKeys(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", id)
		}
		return nil, fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	return flattenFunctionAppFunctionKeys(resp.Model), nil
}

func flattenFunctionAppFunctionKeys(input *webapps.StringDictionary) map[string]string {
	result := make(map[string]string)
	if input == nil || input.Properties == nil {
		return result
	}

	for k, v := range *input.Properties {
		result[k] = v
	}

	return result
}
