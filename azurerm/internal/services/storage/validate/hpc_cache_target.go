package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func HPCCacheTargetName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}
	exp := `^[-0-9a-zA-Z_]{1,31}$`
	p := regexp.MustCompile(exp)
	if !p.MatchString(v) {
		errors = append(errors, fmt.Errorf(`cache target name doesn't comply with regexp: "%s"`, exp))
	}

	return warnings, errors
}

func HPCCacheNamespacePath(i interface{}, k string) (warnings []string, errs []error) {
	return absolutePath(i, k)
}

func HPCCacheNFSExport(i interface{}, k string) (warnings []string, errs []error) {
	return absolutePath(i, k)
}

func HPCCacheNFSTargetPath(i interface{}, k string) (warnings []string, errs []error) {
	return relativePath(i, k)
}

func absolutePath(i interface{}, k string) (warnings []string, errs []error) {
	v, ok := i.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !strings.HasPrefix(v, "/") {
		errs = append(errs, fmt.Errorf(`%s path should start with "/"`, k))
	}
	return warnings, errs
}

func relativePath(i interface{}, k string) (warnings []string, errs []error) {
	v, ok := i.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if strings.HasPrefix(v, "/") {
		errs = append(errs, fmt.Errorf(`%s path should not start with "/"`, k))
	}
	return warnings, errs
}
