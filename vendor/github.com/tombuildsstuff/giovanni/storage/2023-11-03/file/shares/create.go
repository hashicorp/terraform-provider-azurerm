package shares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type AccessTier string

const (
	TransactionOptimizedAccessTier AccessTier = "TransactionOptimized"
	HotAccessTier                  AccessTier = "Hot"
	CoolAccessTier                 AccessTier = "Cool"
	PremiumAccessTier              AccessTier = "Premium"
)

type CreateInput struct {
	// Specifies the maximum size of the share, in gigabytes.
	// Must be greater than 0, and less than or equal to 5TB (5120).
	QuotaInGB int

	// Specifies the enabled protocols on the share. If not specified, the default is SMB.
	EnabledProtocol ShareProtocol

	MetaData map[string]string

	// Specifies the access tier of the share.
	AccessTier *AccessTier
}

type CreateResponse struct {
	HttpResponse *client.Response
}

// Create creates the specified Storage Share within the specified Storage Account
func (c Client) Create(ctx context.Context, shareName string, input CreateInput) (resp CreateResponse, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if input.QuotaInGB <= 0 || input.QuotaInGB > 102400 {
		return resp, fmt.Errorf("`input.QuotaInGB` must be greater than 0, and less than/equal to 100TB (102400 GB)")
	}

	if err = metadata.Validate(input.MetaData); err != nil {
		return resp, fmt.Errorf("`input.MetaData` is not valid: %s", err)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: CreateOptions{
			input: input,
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

type CreateOptions struct {
	input CreateInput
}

func (c CreateOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if len(c.input.MetaData) > 0 {
		headers.Merge(metadata.SetMetaDataHeaders(c.input.MetaData))
	}

	protocol := SMB
	if c.input.EnabledProtocol != "" {
		protocol = c.input.EnabledProtocol
	}
	headers.Append("x-ms-enabled-protocols", string(protocol))

	if c.input.AccessTier != nil {
		headers.Append("x-ms-access-tier", string(*c.input.AccessTier))
	}

	headers.Append("x-ms-share-quota", strconv.Itoa(c.input.QuotaInGB))

	return headers
}

func (c CreateOptions) ToOData() *odata.Query {
	return nil
}

func (c CreateOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	return out
}
