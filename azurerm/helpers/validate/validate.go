package validate

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"net/url"
)

//todo, now in terraform helper, switch over once vended,
func Rfc3339Time(i interface{}, k string) (ws []string, errors []error) {
	if _, err := date.ParseTime(time.RFC3339, i.(string)); err != nil {
		errors = append(errors, fmt.Errorf("%q has the invalid RFC3339 date format %q: %+v", k, i, err))
	}

	return
}

func Url(i interface{}, k string) (ws []string, errors []error) {
	url, err := url.Parse(i.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("%q url is in an invalid format: %q (%+v)", k, i, err))
	} else if url.Scheme != "http" && url.Scheme != "https" {
		errors = append(errors, fmt.Errorf("%q url is neither an http or https uri: %q", k, url.Scheme))
	} else if url.Host == "" {
		errors = append(errors, fmt.Errorf("%q url has no host: %q", k, url))
	}

	return
}
