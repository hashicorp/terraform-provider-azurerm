package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ListRangesResponse struct {
	HttpResponse *client.Response

	Ranges []Range `xml:"Range"`
}

type Range struct {
	Start string `xml:"Start"`
	End   string `xml:"End"`
}

// ListRanges returns the list of valid ranges for the specified File.
func (c Client) ListRanges(ctx context.Context, shareName, path, fileName string) (resp ListRangesResponse, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if path == "" {
		return resp, fmt.Errorf("`path` cannot be an empty string")
	}

	if fileName == "" {
		return resp, fmt.Errorf("`fileName` cannot be an empty string")
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: ListRangeOptions{},
		Path:          fmt.Sprintf("/%s/%s%s", shareName, path, fileName),
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
		resp.HttpResponse.Unmarshal(&resp)
	}

	return
}

type ListRangeOptions struct{}

func (l ListRangeOptions) ToHeaders() *client.Headers {
	return nil
}

func (l ListRangeOptions) ToOData() *odata.Query {
	return nil
}

func (l ListRangeOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "rangelist")
	return out
}
