// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceproviders

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

var ErrNoAuthorization = errors.New("authorization failed")

const registrationErrorV3Fmt = `%s.

Terraform automatically attempts to register the Azure Resource Providers it supports, to
ensure it is able to provision resources.

If you don't have permission to register Resource Providers you may wish to disable this
functionality by adding the following to the Provider block:

provider "azurerm" {
  skip_provider_registration = true
}
Please note that if you opt out of Resource Provider Registration and Terraform tries
to provision a resource from a Resource Provider which is unregistered, then the errors
may appear misleading - for example:
> API version 2019-XX-XX was not found for Microsoft.Foo
Could suggest that the Resource Provider "Microsoft.Foo" requires registration, but
this could also indicate that this Azure Region doesn't support this API version.
More information on the "skip_provider_registration" property can be found here:
https://registry.terraform.io/providers/hashicorp/azurerm/3.116.0/docs#skip_provider_registration
Encountered the following errors:
%v`

const registrationErrorV4Fmt = `%s.
Terraform automatically attempts to register the Azure Resource Providers it supports, to
ensure it is able to provision resources.
If you don't have permission to register Resource Providers you may wish to disable this
functionality by adding the following to the Provider block:

provider "azurerm" {
  resource_provider_registrations = "none"
}

Please note that if you opt out of Resource Provider Registration and Terraform tries
to provision a resource from a Resource Provider which is unregistered, then the errors
may appear misleading - for example:

> API version 2019-XX-XX was not found for Microsoft.Foo

Could suggest that the Resource Provider "Microsoft.Foo" requires registration, but
this could also indicate that this Azure Region doesn't support this API version.

More information on the "resource_provider_registrations" property can be found here:
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#resource-provider-registrations

Encountered the following errors:

%v`

// userError exposes custom messaging depending on the type of error(s) received when attempting to register
// resource providers.
func userError(err error) error {
	if errors.Is(err, ErrNoAuthorization) {
		if !features.FourPointOhBeta() {
			return fmt.Errorf(registrationErrorV3Fmt, "Terraform does not have the necessary permissions to register Resource Providers", err)
		}
		return fmt.Errorf(registrationErrorV4Fmt, "Terraform does not have the necessary permissions to register Resource Providers", err)
	}
	if !features.FourPointOhBeta() {
		return fmt.Errorf(registrationErrorV3Fmt, "Encountered an error whilst ensuring Resource Providers are registered", err)
	}
	return fmt.Errorf(registrationErrorV4Fmt, "Encountered an error whilst ensuring Resource Providers are registered", err)
}

// registrationErrors is a container for errors encountered when attempting to register resource providers. It makes
// use of unwrap support in `errors.Is()` to detect whether any of the contained errors match a target error, which
// allows us to expose different user-facing messages depending on the types of errors encountered. A mutex is
// necessary, as we populate this from goroutines.
type registrationErrors struct {
	errs []error
	lock sync.Mutex
}

func (e *registrationErrors) Error() (out string) {
	return errors.Join(e.errs...).Error()
}

func (e *registrationErrors) Unwrap() []error {
	return e.errs
}

func (e *registrationErrors) append(err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.errs == nil {
		e.errs = make([]error, 0)
	}

	e.errs = append(e.errs, err)
}

func (e *registrationErrors) hasErr() bool {
	return len(e.errs) > 0
}
