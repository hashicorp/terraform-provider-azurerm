package directories

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type SetMetaDataResponse struct {
	HttpResponse *client.Response
}

type SetMetaDataInput struct {
	MetaData map[string]string
}

// SetMetaData updates user defined metadata for the specified directory
func (c Client) SetMetaData(ctx context.Context, shareName, path string, input SetMetaDataInput) (resp SetMetaDataResponse, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if err := metadata.Validate(input.MetaData); err != nil {
		return resp, fmt.Errorf("`metadata` is not valid: %s", err)
	}

	if path == "" {
		return resp, fmt.Errorf("`path` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: SetMetaDataOptions{
			metaData: input.MetaData,
		},
		Path: fmt.Sprintf("/%s/%s", shareName, path),
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

	return
}

type SetMetaDataOptions struct {
	metaData map[string]string
}

func (s SetMetaDataOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if len(s.metaData) > 0 {
		headers.Merge(metadata.SetMetaDataHeaders(s.metaData))
	}
	return headers
}

func (s SetMetaDataOptions) ToOData() *odata.Query {
	return nil
}

func (s SetMetaDataOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "directory")
	out.Append("comp", "metadata")
	return out
}
