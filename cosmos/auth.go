package cosmos

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"net/http"
	"net/url"
	"strings"
)

// Handle adding the correct authorization header to the request
// for more details see: https://docs.microsoft.com/en-us/rest/api/cosmos-db/access-control-on-cosmosdb-resources?redirectedfrom=MSDN
func WithAuthorizationHeader(key, keyType, keyVersion string) autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, fmt.Errorf("Error preparing request: %v", err)
			}

			if r.Header == nil {
				r.Header = make(http.Header)
			}

			signature, err := GenerateAuthorizationSignature(r, key)
			if err != nil {

			}
			r.Header.Set(http.CanonicalHeaderKey("Authorization"), FormatAuthorizationHeader(keyType, keyVersion, signature))

			return r, nil
		})
	}
}

func GenerateAuthorizationSignature(r *http.Request, key64 string) (string, error) {

	// get the required parts
	verb := strings.ToLower(r.Method)
	date := strings.ToLower(r.Header.Get("x-ms-date"))
	rLink := r.URL.Path
	rType := ""

	// normalize rLink to be 'naked'
	rLink = strings.TrimPrefix(rLink, "/")
	rLink = strings.TrimSuffix(rLink, "/")

	// get the type
	// "dbs/MyDatabase/colls/MyCollection" -> colls
	// "dbs" -> colls
	rLinkParts := strings.Split(rLink, "/")
	rLinkPartsLen := len(rLinkParts)

	if rLinkPartsLen%2 == 0 {
		rType = rLinkParts[rLinkPartsLen-2]
	} else {
		rType = rLinkParts[rLinkPartsLen-1]
	}

	//form signature to hash
	sig := verb + "\n" + strings.ToLower(rType) + "\n" + rLink + "\n" + date + "\n\n" //yes the extra new line is required

	//decode key
	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		return "", fmt.Errorf("Error generating authorization signature for %v", err)
	}

	hmac := hmac.New(sha256.New, key)
	hmac.Write([]byte(sig))
	sum := hmac.Sum(nil)

	sum64 := base64.StdEncoding.EncodeToString(sum)
	return sum64, nil
}

func FormatAuthorizationHeader(keyType, keyVersion, signature string) string {
	return url.QueryEscape("type=" + keyType + "&ver=" + keyVersion + "&sig=" + signature)
}
