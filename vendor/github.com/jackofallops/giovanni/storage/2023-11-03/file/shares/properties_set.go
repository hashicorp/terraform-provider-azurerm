package shares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ShareProperties struct {
	QuotaInGb  *int
	AccessTier *AccessTier
}

type SetPropertiesResponse struct {
	HttpResponse *http.Response
}

// SetProperties lets you update the Quota for the specified Storage Share
func (c Client) SetProperties(ctx context.Context, shareName string, properties ShareProperties) (result SetPropertiesResponse, err error) {

	if shareName == "" {
		return result, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return result, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if newQuotaGB := properties.QuotaInGb; newQuotaGB != nil && (*newQuotaGB <= 0 || *newQuotaGB > 102400) {
		return result, fmt.Errorf("`newQuotaGB` must be greater than 0, and less than/equal to 100TB (102400 GB)")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: SetPropertiesOptions{
			input: properties,
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

type SetPropertiesOptions struct {
	input ShareProperties
}

func (s SetPropertiesOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if s.input.QuotaInGb != nil {
		headers.Append("x-ms-share-quota", strconv.Itoa(*s.input.QuotaInGb))
	}

	if s.input.AccessTier != nil {
		headers.Append("x-ms-access-tier", string(*s.input.AccessTier))
	}
	return headers
}

func (s SetPropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (s SetPropertiesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("comp", "properties")
	return out
}
