// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/response"
)

// ResourceCreateRefreshFunc returns a StateRefreshFunc for the standard Azure
// eventual-consistency pattern: poll until Get returns a non-404 response.
func ResourceCreateRefreshFunc(getter func() (*http.Response, error)) StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpResp, err := getter()
		if err != nil {
			if response.WasNotFound(httpResp) {
				return nil, "NotFound", nil
			}
			return nil, "Error", err
		}
		return httpResp, "Found", nil
	}
}
