package apimanagement_test

import (
	"log"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement"
)

func TestXmlWithDotNetInterpolationsDiffSuppress(t *testing.T) {
	testData := []struct {
		old  string
		new  string
		same bool
	}{
		{
			old:  "",
			new:  "",
			same: true,
		},
		{
			old:  "<hello />",
			new:  "",
			same: false,
		},
		{
			old:  "",
			new:  "<hello />",
			same: false,
		},
		{
			old:  "<hello />",
			new:  "<world />",
			same: false,
		},
		{
			// no expressions - with whitespace
			old:  "<policies><inbound><set-variable name=\"abc\" value=\"bcd\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			new:  "<policies><inbound><set-variable name=\"abc\" value=\"bcd\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			same: true,
		},
		{
			// not xml encoded - with whitespace
			old:  "<policies>\n  <inbound>\n    <set-variable name=\"abc\" value=\"bcd\" />\n    <find-and-replace from=\"xyz\" to=\"abc\" />\n  </inbound>\n</policies>\n",
			new:  "<policies>\r\n\t<inbound>\r\n\t\t<set-variable name=\"abc\" value=\"bcd\" />\r\n\t\t<find-and-replace from=\"xyz\" to=\"abc\" />\r\n\t</inbound>\r\n</policies>",
			same: true,
		},
		{
			// both are xml encoded - with whitespace
			old:  "<policies><inbound><set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(\"X-Header-Name\", \"\"))\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			new:  "<policies><inbound><set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(\"X-Header-Name\", \"\"))\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			same: true,
		},
		{
			// not xml encoded - with whitespace
			old:  "<policies>\n  <inbound>\n    <set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(\"X-Header-Name\", \"\"))\" />\n    <find-and-replace from=\"xyz\" to=\"abc\" />\n  </inbound>\n</policies>\n",
			new:  "<policies>\r\n\t<inbound>\r\n\t\t<set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(\"X-Header-Name\", \"\"))\" />\r\n\t\t<find-and-replace from=\"xyz\" to=\"abc\" />\r\n\t</inbound>\r\n</policies>",
			same: true,
		},
		{
			// both are xml encoded - no whitespace
			old:  "<policies><inbound><set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(&quot;X-Header-Name&quot;, &quot;&quot;))\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			new:  "<policies><inbound><set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(&quot;X-Header-Name&quot;, &quot;&quot;))\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			same: true,
		},
		{
			// both are xml encoded with whitespace
			old:  "<policies>\n  <inbound>\n    <set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(&quot;X-Header-Name&quot;, &quot;&quot;))\" />\n    <find-and-replace from=\"xyz\" to=\"abc\" />\n  </inbound>\n</policies>\n",
			new:  "<policies>\r\n\t<inbound>\r\n\t\t<set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(&quot;X-Header-Name&quot;, &quot;&quot;))\" />\r\n\t\t<find-and-replace from=\"xyz\" to=\"abc\" />\r\n\t</inbound>\r\n</policies>",
			same: true,
		},
		{
			// new is xml encoded, old isn't
			old:  "<policies><inbound><set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(\"X-Header-Name\", \"\"))\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			new:  "<policies><inbound><set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(&quot;X-Header-Name&quot;, &quot;&quot;))\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			same: true,
		},
		{
			// new is xml encoded, old isn't
			old:  "<policies>\n  <inbound>\n    <set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(\"X-Header-Name\", \"\"))\" />\n    <find-and-replace from=\"xyz\" to=\"abc\" />\n  </inbound>\n</policies>\n",
			new:  "<policies>\r\n\t<inbound>\r\n\t\t<set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(&quot;X-Header-Name&quot;, &quot;&quot;))\" />\r\n\t\t<find-and-replace from=\"xyz\" to=\"abc\" />\r\n\t</inbound>\r\n</policies>",
			same: true,
		},
		{
			// old is xml encoded, new isn't
			old:  "<policies><inbound><set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(&quot;X-Header-Name&quot;, &quot;&quot;))\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			new:  "<policies><inbound><set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(\"X-Header-Name\", \"\"))\" /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>",
			same: true,
		},
		{
			// old is xml encoded, new isn't
			old:  "<policies>\r\n\t<inbound>\r\n\t\t<set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(&quot;X-Header-Name&quot;, &quot;&quot;))\" />\r\n\t\t<find-and-replace from=\"xyz\" to=\"abc\" />\r\n\t</inbound>\r\n</policies>",
			new:  "<policies>\n  <inbound>\n    <set-variable name=\"abc\" value=\"@(context.Request.Headers.GetValueOrDefault(\"X-Header-Name\", \"\"))\" />\n    <find-and-replace from=\"xyz\" to=\"abc\" />\n  </inbound>\n</policies>\n",
			same: true,
		},
	}

	for _, v := range testData {
		log.Printf("[DEBUG] Testing %q vs %q..", v.old, v.new)
		actual := apimanagement.XmlWithDotNetInterpolationsDiffSuppress("", v.old, v.new, nil)
		if actual != v.same {
			t.Fatalf("Expected %t but got %t", v.same, actual)
		}
	}
}
