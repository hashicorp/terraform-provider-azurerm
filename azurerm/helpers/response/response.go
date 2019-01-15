package response

import (
	"net/http"

	"github.com/hashicorp/go-azure-helpers/response"
)

// TODO: deprecate and remove these
func WasConflict(resp *http.Response) bool {
	return response.WasConflict(resp)
}

// TODO: deprecate and remove these
func WasNotFound(resp *http.Response) bool {
	return response.WasNotFound(resp)
}
