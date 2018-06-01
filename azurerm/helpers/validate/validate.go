package validate

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
)

//todo, now in terraform helper, switch over once vended,
func Rfc3339Time(i interface{}, k string) (_ []string, errors []error) {
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
	return func(i interface{}, k string) (_ []string, errors []error) {
		v, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %q to be int", k))
			return
		}

		if v < min || v > max {
			errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		if v == not {
			errors = append(errors, fmt.Errorf("expected %s to not be %d, got %d", k, not, v))
			return
		}

		return
	}
}

func UrlIsHttpOrHttps() schema.SchemaValidateFunc {
	return UrlWithScheme([]string{"http", "https"})
}

func UrlWithScheme(validSchemes []string) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (_ []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
			return
		}

		url, err := url.Parse(v)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q url is in an invalid format: %q (%+v)", k, i, err))
			return
		}

		if url.Host == "" {
			errors = append(errors, fmt.Errorf("%q url has no host: %q", k, url))
		}

		for _, s := range validSchemes {
			if url.Scheme == s {
				return //last check so just return
			}
		}

		return
	}
}
