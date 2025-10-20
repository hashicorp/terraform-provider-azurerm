package shares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type GetPropertiesResult struct {
	HttpResponse *http.Response

	MetaData        map[string]string
	QuotaInGB       int
	EnabledProtocol ShareProtocol
	AccessTier      *AccessTier
}

// GetProperties returns the properties about the specified Storage Share
func (c Client) GetProperties(ctx context.Context, shareName string) (result GetPropertiesResult, err error) {
	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.MetaData = metadata.ParseFromHeaders(resp.Header)

				quotaRaw := resp.Header.Get("x-ms-share-quota")
				if quotaRaw != "" {
					quota, e := strconv.Atoi(quotaRaw)
					if e != nil {
						err = fmt.Errorf("error converting %q to an integer: %s", quotaRaw, err)
						return
					}
					result.QuotaInGB = quota
				}

				protocol := SMB
				if protocolRaw := resp.Header.Get("x-ms-enabled-protocols"); protocolRaw != "" {
					protocol = ShareProtocol(protocolRaw)
				}

				if accessTierRaw := resp.Header.Get("x-ms-access-tier"); accessTierRaw != "" {
					tier := AccessTier(accessTierRaw)
					result.AccessTier = &tier
				}
				result.EnabledProtocol = protocol
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}
