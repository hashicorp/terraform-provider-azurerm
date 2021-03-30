// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

// Constants ensuring that header names are correctly spelled and consistently cased.
const (
	HeaderAuthorization      = "Authorization"
	HeaderCacheControl       = "Cache-Control"
	HeaderContentEncoding    = "Content-Encoding"
	HeaderContentDisposition = "Content-Disposition"
	HeaderContentLanguage    = "Content-Language"
	HeaderContentLength      = "Content-Length"
	HeaderContentMD5         = "Content-MD5"
	HeaderContentType        = "Content-Type"
	HeaderDate               = "Date"
	HeaderIfMatch            = "If-Match"
	HeaderIfModifiedSince    = "If-Modified-Since"
	HeaderIfNoneMatch        = "If-None-Match"
	HeaderIfUnmodifiedSince  = "If-Unmodified-Since"
	HeaderMetadata           = "Metadata"
	HeaderRange              = "Range"
	HeaderRetryAfter         = "Retry-After"
	HeaderURLEncoded         = "application/x-www-form-urlencoded"
	HeaderUserAgent          = "User-Agent"
	HeaderXmsDate            = "x-ms-date"
	HeaderXmsVersion         = "x-ms-version"
)
