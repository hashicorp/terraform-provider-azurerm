// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package accept

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Header represents an HTTP Accept header value
// See https://httpwg.org/specs/rfc9110.html#field.accept
type Header struct {
	types []PreferredType
}

func (h Header) FirstChoice() *PreferredType {
	if len(h.types) == 0 {
		return nil
	}
	return &h.types[0]
}

func (h Header) String() string {
	out := make([]string, 0)
	for _, typ := range h.types {
		out = append(out, fmt.Sprintf("%s", typ))
	}
	return strings.Join(out, ", ")
}

type PreferredType struct {
	// ContentType is a singular media type (e.g. text/plain), or one partially containing
	// wildcards (e.g. text/*), or entirely wildcards (e.g. */*)
	ContentType string

	// Parameters is a map of media type parameters, e.g. charset=utf-8
	Parameters map[string]string

	// Weight is the integer-normalized quality value representing the relative weight/preference
	Weight qValue
}

func (t PreferredType) String() string {
	out := []string{t.ContentType}
	for p, v := range t.Parameters {
		out = append(out, fmt.Sprintf("%s=%s", p, v))
	}
	out = append(out, fmt.Sprintf("q=%s", t.Weight))

	return strings.Join(out, "; ")
}

// qValue is an integer-normalized representation of a quality value, which has a minimum value
// of 0.001, a maximum value of 1 and a maximum precision of 3 decimal places.
// See https://httpwg.org/specs/rfc9110.html#quality.values
type qValue uint16

func (q qValue) String() string {
	return fmt.Sprintf("%d.%d", q/1000, q%1000)
}

func FromString(in string) (Header, error) {
	contentTypes := strings.Split(in, ",")
	if len(contentTypes) == 0 {
		return Header{}, fmt.Errorf("empty header value provided")
	}

	types := make([]PreferredType, 0)
	for _, typeRaw := range contentTypes {
		// separate the parameters from the media type
		split := strings.Split(strings.ToLower(typeRaw), ";")
		if len(split) == 0 {
			continue
		}

		typ := PreferredType{
			ContentType: strings.TrimSpace(split[0]),
			Parameters:  make(map[string]string),
			Weight:      1000, // default weight of 1.000
		}

		if len(split) > 1 {
			params := split[1:]
			for _, param := range params {
				param = strings.TrimSpace(param)

				if len(param) < 3 || strings.Index(param, "=") < 1 {
					return Header{}, fmt.Errorf("invalid parameter for %q: %q", typ.ContentType, param)
				}

				p := strings.Split(param, "=")
				if len(p) > 2 {
					return Header{}, fmt.Errorf("parameter contains multiple `=` for %q: %q", typ.ContentType, param)
				}

				for i, v := range p {
					p[i] = strings.TrimSpace(v)
				}

				// handle quality values and store other parameters in a map
				switch p[0] {
				case "q":
					q := p[1]

					// determine whether it contains useful decimal point
					pIndex := strings.Index(q, ".")
					if pIndex == len(q) {
						q = q[0 : len(q)-1]
						pIndex = -1
					}

					if pIndex > 0 {
						zIndex := len(q)

						// pad out the value to three decimal places
						if zIndex-pIndex < 3 {
							q = q + strings.Repeat("0", zIndex-pIndex)
							zIndex = len(q)
						}

						// trim any precision > 3 decimal places
						if zIndex-pIndex > 4 {
							zIndex = pIndex + 4
						}

						// strip out the decimal point, so we are left with thousandths
						q = q[0:pIndex] + q[pIndex+1:zIndex]
					}

					// convert the thousandths value to an int
					weight, err := strconv.Atoi(q)
					if err != nil {
						return Header{}, fmt.Errorf("invalid weight for %q: %q", typ.ContentType, q)
					}

					// if an integer was supplied, just multiply it for normalized value
					if pIndex < 1 {
						typ.Weight = qValue(weight * 1000)
					} else {
						typ.Weight = qValue(weight)
					}

				default:
					typ.Parameters[p[0]] = p[1]
				}
			}
		}

		types = append(types, typ)
	}

	sort.Slice(types, func(i, j int) bool {
		return types[i].Weight > types[j].Weight
	})

	return Header{types: types}, nil
}
