// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var _ Authorizer = &SharedKeyAuthorizer{}

// SharedKeyType defines the enumeration for the various shared key types.
// See https://docs.microsoft.com/en-us/rest/api/storageservices/authorize-with-shared-key for details on the shared key types.
type SharedKeyType string

const (
	SharedKey      SharedKeyType = "sharedKey"
	SharedKeyTable SharedKeyType = "sharedKeyTable"
)

type SharedKeyAuthorizer struct {
	accountName string
	accountKey  []byte
	keyType     SharedKeyType
}

func NewSharedKeyAuthorizer(accountName string, accountKey string, keyType SharedKeyType) (*SharedKeyAuthorizer, error) {
	key, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return nil, fmt.Errorf("decoding accountKey: %+v", err)
	}
	return &SharedKeyAuthorizer{
		accountName: accountName,
		accountKey:  key,
		keyType:     keyType,
	}, nil
}

func (s *SharedKeyAuthorizer) Token(ctx context.Context, req *http.Request) (*oauth2.Token, error) {
	key, err := buildSharedKey(s.accountName, s.accountKey, req, s.keyType)
	if err != nil {
		return nil, fmt.Errorf("building SharedKey for request: %+v", err)
	}
	return &oauth2.Token{
		TokenType:   "SharedKey",
		AccessToken: key,
	}, nil
}

func (s *SharedKeyAuthorizer) AuxiliaryTokens(_ context.Context, _ *http.Request) ([]*oauth2.Token, error) {
	// Auxiliary tokens are not supported with SharedKey authentication
	return []*oauth2.Token{}, nil
}

const (
	storageEmulatorAccountName string = "devstoreaccount1"

	headerContentEncoding   = "Content-Encoding"
	headerContentMD5        = "Content-MD5"
	headerContentLanguage   = "Content-Language"
	headerContentType       = "Content-Type"
	headerIfModifiedSince   = "If-Modified-Since"
	headerIfMatch           = "If-Match"
	headerIfNoneMatch       = "If-None-Match"
	headerIfUnmodifiedSince = "If-Unmodified-Since"
	headerDate              = "Date"
	headerXMSDate           = "X-Ms-Date"
	headerRange             = "Range"
)

func buildSharedKey(accName string, accKey []byte, req *http.Request, keyType SharedKeyType) (string, error) {
	canRes, err := buildCanonicalizedResource(accName, req.URL.String(), keyType)
	if err != nil {
		return "", err
	}

	if req.Header == nil {
		req.Header = http.Header{}
	}

	// ensure date is set
	if req.Header.Get(headerDate) == "" && req.Header.Get(headerXMSDate) == "" {
		date := time.Now().UTC().Format(http.TimeFormat)
		req.Header.Set(headerXMSDate, date)
	}
	canString, err := buildCanonicalizedString(req.Method, req.ContentLength, req.Header, canRes, keyType)
	if err != nil {
		return "", err
	}
	return createAuthorizationHeader(accName, accKey, canString), nil
}

func buildCanonicalizedResource(accountName, uri string, keyType SharedKeyType) (string, error) {
	errMsg := "buildCanonicalizedResource error: %s"
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf(errMsg, err.Error())
	}

	cr := bytes.NewBufferString("")
	if accountName != storageEmulatorAccountName {
		cr.WriteString("/")
		cr.WriteString(getCanonicalizedAccountName(accountName))
	}

	if len(u.Path) > 0 {
		// Any portion of the CanonicalizedResource string that is derived from
		// the resource's URI should be encoded exactly as it is in the URI.
		// -- https://msdn.microsoft.com/en-gb/library/azure/dd179428.aspx
		cr.WriteString(u.EscapedPath())
	} else {
		// a slash is required to indicate the root path
		cr.WriteString("/")
	}

	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", fmt.Errorf(errMsg, err.Error())
	}

	// See https://github.com/Azure/azure-storage-net/blob/master/Lib/Common/Core/Util/AuthenticationUtility.cs#L277
	if keyType == SharedKey {
		if len(params) > 0 {
			cr.WriteString("\n")

			keys := []string{}
			for key := range params {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			completeParams := []string{}
			for _, key := range keys {
				if len(params[key]) > 1 {
					sort.Strings(params[key])
				}

				completeParams = append(completeParams, fmt.Sprintf("%s:%s", key, strings.Join(params[key], ",")))
			}
			cr.WriteString(strings.Join(completeParams, "\n"))
		}
	} else {
		// search for "comp" parameter, if exists then add it to canonicalizedresource
		if v, ok := params["comp"]; ok {
			cr.WriteString("?comp=" + v[0])
		}
	}

	return cr.String(), nil
}

func getCanonicalizedAccountName(accountName string) string {
	// since we may be trying to access a secondary storage account, we need to
	// remove the -secondary part of the storage name
	return strings.TrimSuffix(accountName, "-secondary")
}

func buildCanonicalizedString(verb string, contentLength int64, headers http.Header, canonicalizedResource string, keyType SharedKeyType) (string, error) {
	var contentLengthString string
	if contentLength > 0 {
		contentLengthString = strconv.Itoa(int(contentLength))
	}
	date := headers.Get(headerDate)
	if v := headers.Get(headerXMSDate); v != "" {
		if keyType == SharedKey {
			date = ""
		} else {
			date = v
		}
	}
	var canString string
	switch keyType {
	case SharedKey:
		canString = strings.Join([]string{
			verb,
			headers.Get(headerContentEncoding),
			headers.Get(headerContentLanguage),
			contentLengthString,
			headers.Get(headerContentMD5),
			headers.Get(headerContentType),
			date,
			headers.Get(headerIfModifiedSince),
			headers.Get(headerIfMatch),
			headers.Get(headerIfNoneMatch),
			headers.Get(headerIfUnmodifiedSince),
			headers.Get(headerRange),
			buildCanonicalizedHeader(headers),
			canonicalizedResource,
		}, "\n")
	case SharedKeyTable:
		canString = strings.Join([]string{
			verb,
			headers.Get(headerContentMD5),
			headers.Get(headerContentType),
			date,
			canonicalizedResource,
		}, "\n")
	default:
		return "", fmt.Errorf("key type '%s' is not supported", keyType)
	}
	return canString, nil
}

func buildCanonicalizedHeader(headers http.Header) string {
	cm := make(map[string]string)

	for k := range headers {
		headerName := strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(headerName, "x-ms-") {
			cm[headerName] = headers.Get(k)
		}
	}

	if len(cm) == 0 {
		return ""
	}

	keys := []string{}
	for key := range cm {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	ch := bytes.NewBufferString("")

	for _, key := range keys {
		ch.WriteString(key)
		ch.WriteRune(':')
		ch.WriteString(cm[key])
		ch.WriteRune('\n')
	}

	return strings.TrimSuffix(ch.String(), "\n")
}

func createAuthorizationHeader(accountName string, accountKey []byte, canonicalizedString string) string {
	h := hmac.New(sha256.New, accountKey)
	h.Write([]byte(canonicalizedString))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s:%s", getCanonicalizedAccountName(accountName), signature)
}
