package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"regexp"
)

//^iqn\.\d{4}-\d{2}.(.[a-zA-Z0-9\-]+){1,63}.(:[\S]+)?$
// panic: regexp: Compile(`^iqn\.\d{4}-\d{2}((?<!-)\.(?!-)[a-zA-Z0-9\-]+){1,63}(?<!-)(?<!\.)(:[\S]+)?$`): error parsing regexp: invalid or unsupported Perl syntax: `(?<`
// IQN A regex converted from https://github.com/rhinstaller/anaconda/blob/master/pyanaconda/core/regexes.py#L163 perl syntax to golang
// a. Starts with string 'iqn.'
// #    b. A date code specifying the year and month in which the organization
// #       registered the domain or sub-domain name used as the naming authority
// #       string. "yyyy-mm"
// #    c. A dot (".")
// #    d. The organizational naming authority string, which consists of a
// #       valid, reversed domain or subdomain name.
// #    e. Optionally, a colon (":"), followed by a string of the assigning
// #       organization's choosing, which must make each assigned iSCSI name
// #       unique. With the exception of the colon prefix, the owner of the domain
// #       name can assign everything after the reversed domain name as desired.

// IQN should follow the format `iqn.yyyy-mm.<abc>.<pqr>[:xyz]`; supported characters include [0-9a-z-.:]
// TODO: wait for regex from svc team
var IQN = validation.StringMatch( //TODO: add unit test, re-write regex
	regexp.MustCompile(`^iqn\.\d{4}-\d{2}.(.[a-zA-Z\d\-]+){1,63}:?[\da-z-.:]$`),
	"IQN should follow the format `iqn.yyyy-mm.<abc>.<pqr>[:xyz]`; supported characters include [0-9a-z-.:]",
	)
