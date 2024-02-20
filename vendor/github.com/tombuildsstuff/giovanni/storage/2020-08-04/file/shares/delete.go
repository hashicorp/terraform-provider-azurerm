package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type DeleteResponse struct {
	HttpResponse *client.Response
}

type DeleteInput struct {
	DeleteSnapshots bool
}

// Delete deletes the specified Storage Share from within a Storage Account
func (c Client) Delete(ctx context.Context, shareName string, input DeleteInput) (resp DeleteResponse, err error) {
	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodDelete,
		OptionsObject: DeleteOptions{
			deleteSnapshots: input.DeleteSnapshots,
		},
		Path: fmt.Sprintf("/%s", shareName),
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

type DeleteOptions struct {
	deleteSnapshots bool
}

func (d DeleteOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if d.deleteSnapshots {
		headers.Append("x-ms-delete-snapshots", "include")
	}
	return headers
}

func (d DeleteOptions) ToOData() *odata.Query {
	return nil
}

func (d DeleteOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	return out
}
