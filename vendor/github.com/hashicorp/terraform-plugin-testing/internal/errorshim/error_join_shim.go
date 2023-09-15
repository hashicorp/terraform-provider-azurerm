// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// TODO: Once Go 1.20 is the minimum supported version delete this package, replace all usages with `errors` package
// - https://github.com/hashicorp/terraform-plugin-testing/issues/99
package errorshim

// Copied from -> https://cs.opensource.google/go/go/+/refs/tags/go1.20.2:src/errors/join.go
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

func (e *joinError) Unwrap() []error {
	return e.errs
}
