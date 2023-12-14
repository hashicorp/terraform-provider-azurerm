package filesystems

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type GetPropertiesResponse struct {
	HttpResponse *client.Response

	// A map of base64-encoded strings to store as user-defined properties with the File System
	// Note that items may only contain ASCII characters in the ISO-8859-1 character set.
	// This automatically gets converted to a comma-separated list of name and
	// value pairs before sending to the API
	Properties map[string]string

	// Is Hierarchical Namespace Enabled?
	NamespaceEnabled bool
}

// GetProperties gets the properties for a Data Lake Store Gen2 FileSystem within a Storage Account
func (c Client) GetProperties(ctx context.Context, fileSystemName string) (resp GetPropertiesResponse, err error) {

	if fileSystemName == "" {
		return resp, fmt.Errorf("`fileSystemName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodHead,
		OptionsObject: fileSystemOptions{},
		Path:          fmt.Sprintf("/%s", fileSystemName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {

		propertiesRaw := resp.HttpResponse.Header.Get("x-ms-properties")
		var properties *map[string]string
		properties, err = parseProperties(propertiesRaw)
		if err != nil {
			return
		}

		resp.Properties = *properties
		resp.NamespaceEnabled = strings.EqualFold(resp.HttpResponse.Header.Get("x-ms-namespace-enabled"), "true")

	}
	return
}
