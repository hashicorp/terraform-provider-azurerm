package validate

import "testing"

func TestDomainName(t *testing.T) {
	cases := []struct {
		Name  string
		Input string
		Valid bool
	}{
		{
			Name:  "random string",
			Input: "dfdsdfds",
			Valid: false,
		},
		{
			Name:  "contains protocol scheme",
			Input: "https://contoso.com",
			Valid: false,
		},
		{
			Name:  "too long",
			Input: "this.hostname.is.definitely.going.to.be.altogether.far.too.long.for.a.valid.rfc.compatible.hostname.even.if.i.have.to.add.a.ludicrous.number.of.parts.to.this.test.case.input.string.including.some.random.character.strings.no.really.i.will.do.it.contoso.com",
			Valid: false,
		},
		{
			Name:  "valid",
			Input: "mydomain.contoso.com",
			Valid: true,
		},
		{
			Name:  "subdomain valid",
			Input: "subdomain.mydomain.contoso.com",
			Valid: true,
		},
	}
	for _, tc := range cases {
		_, err := DomainName(tc.Input, "")

		valid := err == nil
		if valid != tc.Valid {
			t.Errorf("Expected valid status %t but got %t for input %s: %+v", tc.Valid, valid, tc.Input, err)
		}
	}
}
