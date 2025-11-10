package directories

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type GetResponse struct {
	HttpResponse *http.Response

	// A set of name-value pairs that contain metadata for the directory.
	MetaData map[string]string

	// The value of this header is set to true if the directory metadata is completely
	// encrypted using the specified algorithm. Otherwise, the value is set to false.
	DirectoryMetaDataEncrypted bool
}

// Get returns all system properties for the specified directory,
// and can also be used to check the existence of a directory.
func (c Client) Get(ctx context.Context, shareName, path string) (result GetResponse, err error) {
	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if path == "" {
		err = fmt.Errorf("`path` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: directoriesOptions{},
		Path:          fmt.Sprintf("/%s/%s", shareName, path),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.DirectoryMetaDataEncrypted = strings.EqualFold(resp.Header.Get("x-ms-server-encrypted"), "true")
				result.MetaData = metadata.ParseFromHeaders(resp.Header)
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}
