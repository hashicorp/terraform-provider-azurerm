package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type DeleteSnapshotResponse struct {
	HttpResponse *http.Response
}

// DeleteSnapshot deletes the specified Snapshot of a Storage Share
func (c Client) DeleteSnapshot(ctx context.Context, accountName, shareName string, shareSnapshot string) (result DeleteSnapshotResponse, err error) {

	if shareName == "" {
		return result, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return result, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if shareSnapshot == "" {
		return result, fmt.Errorf("`shareSnapshot` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodDelete,
		OptionsObject: snapShotDeleteOptions{
			shareSnapShot: shareSnapshot,
		},
		Path: fmt.Sprintf("/%s", shareName),
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
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type snapShotDeleteOptions struct {
	shareSnapShot string
}

func (s snapShotDeleteOptions) ToHeaders() *client.Headers {
	return nil
}

func (s snapShotDeleteOptions) ToOData() *odata.Query {
	return nil
}

func (s snapShotDeleteOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("sharesnapshot", s.shareSnapShot)
	return out
}
