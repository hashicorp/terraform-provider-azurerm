package files

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetResponse struct {
	HttpResponse *client.Response

	CacheControl          string
	ContentDisposition    string
	ContentEncoding       string
	ContentLanguage       string
	ContentLength         *int64
	ContentMD5            string
	ContentType           string
	CopyID                string
	CopyStatus            string
	CopySource            string
	CopyProgress          string
	CopyStatusDescription string
	CopyCompletionTime    string
	Encrypted             bool

	MetaData map[string]string
}

// GetProperties returns the Properties for the specified file
func (c Client) GetProperties(ctx context.Context, shareName, path, fileName string) (resp GetResponse, err error) {
	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if fileName == "" {
		return resp, fmt.Errorf("`fileName` cannot be an empty string")
	}

	if path != "" {
		path = fmt.Sprintf("/%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodHead,
		OptionsObject: nil,
		Path:          fmt.Sprintf("%s/%s%s", shareName, path, fileName),
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
			resp.CacheControl = resp.HttpResponse.Header.Get("Cache-Control")
			resp.ContentDisposition = resp.HttpResponse.Header.Get("Content-Disposition")
			resp.ContentEncoding = resp.HttpResponse.Header.Get("Content-Encoding")
			resp.ContentLanguage = resp.HttpResponse.Header.Get("Content-Language")
			resp.ContentMD5 = resp.HttpResponse.Header.Get("Content-MD5")
			resp.ContentType = resp.HttpResponse.Header.Get("Content-Type")
			resp.CopyID = resp.HttpResponse.Header.Get("x-ms-copy-id")
			resp.CopyProgress = resp.HttpResponse.Header.Get("x-ms-copy-progress")
			resp.CopySource = resp.HttpResponse.Header.Get("x-ms-copy-source")
			resp.CopyStatus = resp.HttpResponse.Header.Get("x-ms-copy-status")
			resp.CopyStatusDescription = resp.HttpResponse.Header.Get("x-ms-copy-status-description")
			resp.CopyCompletionTime = resp.HttpResponse.Header.Get("x-ms-copy-completion-time")
			resp.Encrypted = strings.EqualFold(resp.HttpResponse.Header.Get("x-ms-server-encrypted"), "true")
			resp.MetaData = metadata.ParseFromHeaders(resp.HttpResponse.Header)

			contentLengthRaw := resp.HttpResponse.Header.Get("Content-Length")
			if contentLengthRaw != "" {
				contentLength, err := strconv.Atoi(contentLengthRaw)
				if err != nil {
					return resp, fmt.Errorf("error parsing %q for Content-Length as an integer: %s", contentLengthRaw, err)
				}
				contentLengthI64 := int64(contentLength)
				resp.ContentLength = &contentLengthI64
			}
		}
	}

	return
}
