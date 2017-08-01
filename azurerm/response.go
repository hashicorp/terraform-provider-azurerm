package azurerm

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

func responseWasNotFound(resp autorest.Response) bool {
	if r := resp.Response; r != nil {
		if r.StatusCode == http.StatusNotFound {
			return true
		}
	}

	return false
}
