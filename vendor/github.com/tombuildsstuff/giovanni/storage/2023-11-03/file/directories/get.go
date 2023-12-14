package directories

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetResponse struct {
	HttpResponse *client.Response

	// A set of name-value pairs that contain metadata for the directory.
	MetaData map[string]string

	// The value of this header is set to true if the directory metadata is completely
	// encrypted using the specified algorithm. Otherwise, the value is set to false.
	DirectoryMetaDataEncrypted bool
}

// Get returns all system properties for the specified directory,
// and can also be used to check the existence of a directory.
func (c Client) Get(ctx context.Context, shareName, path string) (resp GetResponse, err error) {
	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if path == "" {
		return resp, fmt.Errorf("`path` cannot be an empty string")
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

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		if resp.HttpResponse.Header != nil {
			resp.MetaData = metadata.ParseFromHeaders(resp.HttpResponse.Header)
		}
		resp.DirectoryMetaDataEncrypted = strings.EqualFold(resp.HttpResponse.Header.Get("x-ms-server-encrypted"), "true")
	}

	return
}
