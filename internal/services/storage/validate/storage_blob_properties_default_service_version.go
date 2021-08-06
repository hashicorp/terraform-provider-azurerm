package validate

import (
	"fmt"
)

func BlobPropertiesDefaultServiceVersion(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	valid := []string{
		"2008-10-27",
		"2009-04-14",
		"2009-07-17",
		"2009-09-19",
		"2011-08-28",
		"2012-02-12",
		"2013-08-15",
		"2014-02-14",
		"2015-02-21",
		"2015-04-05",
		"2015-07-08",
		"2015-12-11",
		"2016-05-31",
		"2017-04-17",
		"2017-07-29",
		"2017-11-09",
		"2018-03-28",
		"2018-11-09",
		"2019-02-02",
		"2019-07-07",
		"2019-12-12",
		"2020-02-10",
		"2020-04-08",
		"2020-06-12",
	}
	for _, str := range valid {
		if v == str {
			return warnings, errors
		}
	}

	errors = append(errors, fmt.Errorf("expected %s to be one of %v, got %s", k, valid, v))
	return warnings, errors
}
