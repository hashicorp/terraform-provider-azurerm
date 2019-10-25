package suppress

import "testing"

func TestXmlDiff(t *testing.T) {
	cases := []struct {
		Name     string
		XmlA     string
		XmlB     string
		Suppress bool
	}{
		{
			Name:     "empty",
			XmlA:     "",
			XmlB:     "",
			Suppress: true,
		},
		{
			Name:     "neither are xml",
			XmlA:     "this is not an xml",
			XmlB:     "neither is this",
			Suppress: false,
		},
		{
			Name:     "identical texts",
			XmlA:     "this is not an xml",
			XmlB:     "this is not an xml",
			Suppress: true,
		},
		{
			Name:     "xml vs text",
			XmlA:     "<r></r>",
			XmlB:     "this is not an xml",
			Suppress: false,
		},
		{
			Name:     "text vs xml",
			XmlA:     "this is not an xml",
			XmlB:     "<r></r>",
			Suppress: false,
		},
		{
			Name:     "identical xml",
			XmlA:     "<r><c></c></r>",
			XmlB:     "<r><c></c></r>",
			Suppress: true,
		},
		{
			Name:     "xml with different line endings",
			XmlA:     "<r>\n<c>\n</c>\n</r>",
			XmlB:     "<r>\r\n<c>\r\n</c>\r\n</r>",
			Suppress: true,
		},
		{
			Name:     "xml with different indentations",
			XmlA:     "<r>\n  <c>\n  </c>\n</r>",
			XmlB:     "<r>\r\n\t<c>\r\n\t</c>\r\n</r>",
			Suppress: true,
		},
		{
			Name:     "xml with different quotation marks",
			XmlA:     "<r><c attr=\"test\"></c></r>",
			XmlB:     "<r>\r\n\t<c attr='test'>\r\n\t</c>\r\n</r>",
			Suppress: true,
		},
		{
			Name:     "xml with different spaces",
			XmlA:     "<r><c   attr = 'test'></c></r>",
			XmlB:     "<r>\r\n\t<c attr='test'>\r\n\t</c>\r\n</r>",
			Suppress: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if XmlDiff("test", tc.XmlA, tc.XmlB, nil) != tc.Suppress {
				t.Fatalf("Expected XmlDiff to return %t for '%q' == '%q'", tc.Suppress, tc.XmlA, tc.XmlB)
			}
		})
	}
}
