package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetStatsResponse struct {
	HttpResponse *client.Response

	// The approximate size of the data stored on the share.
	// Note that this value may not include all recently created or recently resized files.
	ShareUsageBytes int64 `xml:"ShareUsageBytes"`
}

// GetStats returns information about the specified Storage Share
func (c Client) GetStats(ctx context.Context, shareName string) (resp GetStatsResponse, err error) {
	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: statsOptions{},
		Path:          shareName,
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
		err = resp.HttpResponse.Unmarshal(&resp)
		if err != nil {
			return resp, fmt.Errorf("unmarshalling response: %v", err)
		}
	}
	return
}

type statsOptions struct{}

func (s statsOptions) ToHeaders() *client.Headers {
	return nil
}

func (s statsOptions) ToOData() *odata.Query {
	return nil
}

func (s statsOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("comp", "stats")
	return out
}
