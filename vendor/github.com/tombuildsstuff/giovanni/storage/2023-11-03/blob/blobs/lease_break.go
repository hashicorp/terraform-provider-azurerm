package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type BreakLeaseInput struct {
	//  For a break operation, proposed duration the lease should continue
	//  before it is broken, in seconds, between 0 and 60.
	//  This break period is only used if it is shorter than the time remaining on the lease.
	//  If longer, the time remaining on the lease is used.
	//  A new lease will not be available before the break period has expired,
	//  but the lease may be held for longer than the break period.
	//  If this header does not appear with a break operation, a fixed-duration lease breaks
	//  after the remaining lease period elapses, and an infinite lease breaks immediately.
	BreakPeriod *int

	LeaseID string
}

type BreakLeaseResponse struct {
	HttpResponse *client.Response

	// Approximate time remaining in the lease period, in seconds.
	// If the break is immediate, 0 is returned.
	LeaseTime int
}

// BreakLease breaks an existing lock on a blob using the LeaseID.
func (c Client) BreakLease(ctx context.Context, containerName, blobName string, input BreakLeaseInput) (resp BreakLeaseResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.LeaseID == "" {
		return resp, fmt.Errorf("`input.LeaseID` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: breakLeaseOptions{
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

	return
}

type breakLeaseOptions struct {
	input BreakLeaseInput
}

func (b breakLeaseOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	headers.Append("x-ms-lease-action", "break")
	headers.Append("x-ms-lease-id", b.input.LeaseID)

	if b.input.BreakPeriod != nil {
		headers.Append("x-ms-lease-break-period", strconv.Itoa(*b.input.BreakPeriod))
	}

	return headers
}

func (b breakLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (b breakLeaseOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "lease")
	return out
}
