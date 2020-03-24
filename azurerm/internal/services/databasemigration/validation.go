package databasemigration

import (
	"fmt"
	"regexp"
)

func validateDatabaseMigrationServiceName(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	validName := regexp.MustCompile(`^[\d\w]+[\d\w\-_.]*$`)
	if !validName.MatchString(v) {
		return nil, []error{fmt.Errorf("invalid format of %q", k)}
	}
	return nil, nil
}

func validateDatabaseMigrationProjectName(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	validName := regexp.MustCompile(`^[\d\w]+[\d\w\-_.]*$`)
	if !validName.MatchString(v) {
		return nil, []error{fmt.Errorf("invalid format of %q", k)}
	}
	return nil, nil
}
