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
	HttpResponse *client.Response

	LeaseID string
}

// ChangeLease changes an existing lock on a blob for another lock.
func (c Client) ChangeLease(ctx context.Context, containerName, blobName string, input ChangeLeaseInput) (resp ChangeLeaseResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.ExistingLeaseID == "" {
		return resp, fmt.Errorf("`input.ExistingLeaseID` cannot be an empty string")
	}

	if input.ProposedLeaseID == "" {
		return resp, fmt.Errorf("`input.ProposedLeaseID` cannot be an empty string")
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

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		if resp.HttpResponse.Header != nil {
			resp.LeaseID = resp.HttpResponse.Header.Get("x-ms-lease-id")
		}
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
