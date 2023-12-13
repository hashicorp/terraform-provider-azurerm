package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type SetTierInput struct {
	Tier AccessTier
}

type SetTierResponse struct {
	HttpResponse *client.Response
}

// SetTier sets the tier on a blob.
func (c Client) SetTier(ctx context.Context, containerName, blobName string, input SetTierInput) (resp SetTierResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: setTierOptions{
			tier: input.Tier,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
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

type setTierOptions struct {
	tier AccessTier
}

func (s setTierOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-access-tier", string(s.tier))
	return headers
}

func (s setTierOptions) ToOData() *odata.Query {
	return nil
}

func (s setTierOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "tier")
	return out
}
