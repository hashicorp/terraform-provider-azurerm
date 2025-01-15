package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ChangeLeaseInput struct {
	ExistingLeaseID string
	ProposedLeaseID string
}

type ChangeLeaseResponse struct {
	ChangeLeaseModel
	HttpResponse *http.Response
}

type ChangeLeaseModel struct {
	LeaseID string
}

// ChangeLease changes the lock from one Lease ID to another Lease ID
func (c Client) ChangeLease(ctx context.Context, containerName string, input ChangeLeaseInput) (result ChangeLeaseResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
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
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: changeLeaseOptions{
			existingLeaseId: input.ExistingLeaseID,
			proposedLeaseId: input.ProposedLeaseID,
		},
		Path: fmt.Sprintf("/%s", containerName),
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

var _ client.Options = changeLeaseOptions{}

type changeLeaseOptions struct {
	existingLeaseId string
	proposedLeaseId string
}

func (o changeLeaseOptions) ToHeaders() *client.Headers {
	headers := containerOptions{}.ToHeaders()

	headers.Append("x-ms-lease-action", "change")
	headers.Append("x-ms-lease-id", o.existingLeaseId)
	headers.Append("x-ms-proposed-lease-id", o.proposedLeaseId)

	return headers
}

func (o changeLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (o changeLeaseOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "lease")
	return query
}
