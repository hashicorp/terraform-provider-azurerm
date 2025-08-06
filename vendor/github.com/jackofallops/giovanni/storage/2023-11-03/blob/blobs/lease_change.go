package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ChangeLeaseInput struct {
	ExistingLeaseID string
	ProposedLeaseID string
}

type ChangeLeaseResponse struct {
	HttpResponse *http.Response

	LeaseID string
}

// ChangeLease changes an existing lock on a blob for another lock.
func (c Client) ChangeLease(ctx context.Context, containerName, blobName string, input ChangeLeaseInput) (result ChangeLeaseResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}

	if strings.ToLower(containerName) != containerName {
		err = fmt.Errorf("`containerName` must be a lower-cased string")
		return
	}

	if blobName == "" {
		err = fmt.Errorf("`blobName` cannot be an empty string")
		return
	}

	if input.ExistingLeaseID == "" {
		err = fmt.Errorf("`input.ExistingLeaseID` cannot be an empty string")
		return
	}

	if input.ProposedLeaseID == "" {
		err = fmt.Errorf("`input.ProposedLeaseID` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: changeLeaseOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
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
			if resp.Response != nil && resp.Header != nil {
				result.LeaseID = resp.Header.Get("x-ms-lease-id")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type changeLeaseOptions struct {
	input ChangeLeaseInput
}

func (c changeLeaseOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-lease-action", "change")
	headers.Append("x-ms-lease-id", c.input.ExistingLeaseID)
	headers.Append("x-ms-proposed-lease-id", c.input.ProposedLeaseID)
	return headers
}

func (c changeLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (c changeLeaseOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "lease")
	return out
}
