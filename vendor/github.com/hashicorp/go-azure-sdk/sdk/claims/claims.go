// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package claims

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/oauth2"
)

// Claims is used to unmarshall the claims from a JWT issued by the Microsoft Identity Platform.
type Claims struct {
	Audience          string   `json:"aud"`
	IssuedAt          int64    `json:"iat"`
	Issuer            string   `json:"iss"`
	IdentityProvider  string   `json:"idp"`
	ObjectId          string   `json:"oid"`
	Roles             []string `json:"roles"`
	Scopes            string   `json:"scp"`
	Subject           string   `json:"sub"`
	TenantRegionScope string   `json:"tenant_region_scope"`
	TenantId          string   `json:"tid"`
	Version           string   `json:"ver"`

	AppDisplayName string `json:"app_displayname,omitempty"`
	AppId          string `json:"appid,omitempty"`
	IdType         string `json:"idtyp,omitempty"`
}

// ParseClaims retrieves and parses the claims from a JWT issued by the Microsoft Identity Platform.
func ParseClaims(token *oauth2.Token) (*Claims, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	}

	jwt := strings.Split(token.AccessToken, ".")
	if len(jwt) != 3 {
		return nil, errors.New("unexpected token format: does not have 3 parts")
	}

	payload, err := base64.RawURLEncoding.DecodeString(jwt[1])
	if err != nil {
		return nil, err
	}

	var claims Claims
	err = json.Unmarshal(payload, &claims)
	return &claims, err
}
