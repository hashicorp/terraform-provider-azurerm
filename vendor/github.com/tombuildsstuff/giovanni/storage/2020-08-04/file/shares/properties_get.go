package shares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetPropertiesResult struct {
	HttpResponse *client.Response

	MetaData        map[string]string
	QuotaInGB       int
	EnabledProtocol ShareProtocol
	AccessTier      *AccessTier
}

// GetProperties returns the properties about the specified Storage Share
func (c Client) GetProperties(ctx context.Context, shareName string) (resp GetPropertiesResult, err error) {
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
		OptionsObject: sharesOptions{},
		Path:          fmt.Sprintf("/%s", shareName),
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
			resp.MetaData = metadata.ParseFromHeaders(resp.HttpResponse.Header)

			quotaRaw := resp.HttpResponse.Header.Get("x-ms-share-quota")
			if quotaRaw != "" {
				quota, e := strconv.Atoi(quotaRaw)
				if e != nil {
					return resp, fmt.Errorf("error converting %q to an integer: %s", quotaRaw, err)
				}
				resp.QuotaInGB = quota
			}

			protocol := SMB
			if protocolRaw := resp.HttpResponse.Header.Get("x-ms-enabled-protocols"); protocolRaw != "" {
				protocol = ShareProtocol(protocolRaw)
			}

			if accessTierRaw := resp.HttpResponse.Header.Get("x-ms-access-tier"); accessTierRaw != "" {
				tier := AccessTier(accessTierRaw)
				resp.AccessTier = &tier
			}
			resp.EnabledProtocol = protocol
		}
	}

	return
}
