package validate

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
)

//todo, now in terraform helper, switch over once vended,
func Rfc3339Time(i interface{}, k string) (ws []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := date.ParseTime(time.RFC3339, v); err != nil {
		errors = append(errors, fmt.Errorf("%q has the invalid RFC3339 date format %q: %+v", k, i, err))
	}

	return
}

func IntBetweenAndNot(min, max, not int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %q to be int", k))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		if v == not {
			es = append(es, fmt.Errorf("expected %s to not be %d, got %d", k, not, v))
			return
		}

		return
	}
}

func UrlIsHttp(i interface{}, k string) (ws []string, errors []error) {

}

func Url(i interface{}, k string) (ws []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	url, err := url.Parse(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q url is in an invalid format: %q (%+v)", k, i, err))
	} else if url.Scheme != "http" && url.Scheme != "https" {
		errors = append(errors, fmt.Errorf("%q url is neither an http or https uri: %q", k, url.Scheme))
	} else if url.Host == "" {
		errors = append(errors, fmt.Errorf("%q url has no host: %q", k, url))
	}

	return
}
