package azurerm

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/net/html"
)

type Result struct {
	XMLName xml.Name `xml:"AzurePublicIpAddresses"`
	Regions []Region `xml:"Region"`
}

type Region struct {
	Name     string    `xml:"Name,attr"`
	IpRanges []IpRange `xml:"IpRange"`
}

type IpRange struct {
	Subnet string `xml:"Subnet,attr"`
}

func dataSourceIpRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIpRangesRead,

		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIpRangesRead(d *schema.ResourceData, meta interface{}) error {
	regionsSet := d.Get("regions").(*schema.Set)

	url := "https://www.microsoft.com/en-us/download/confirmation.aspx?id=41653"
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error performing GET on url (%s): %s", url, err)
	}

	xmlUrl, err := getLinkToPublicIpsXML(resp.Body)
	if err != nil {
		return fmt.Errorf("Error extracting link to XML document: %s", err)
	}
	d.SetId(xmlUrl)

	log.Printf("[DEBUG] Reading IP ranges from %s", xmlUrl)

	resp, err = http.Get(xmlUrl)
	if err != nil {
		return fmt.Errorf("Error extracting XML from url (%s): %s", url, err)
	}

	result := &Result{}
	err = xml.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("Error decoding XML: %s", err)
	}

	var subnets []string
	for _, region := range result.Regions {
		if regionsSet.Len() == 0 ||
			regionsSet.Contains(strings.ToLower(region.Name)) {
			for _, ipRange := range region.IpRanges {
				subnets = append(subnets, ipRange.Subnet)
			}
		}
	}

	if len(subnets) == 0 {
		return fmt.Errorf("No IP ranges results from regions: %s",
			regionsSet.GoString(),
		)
	}

	sort.Strings(subnets)
	if err := d.Set("subnets", subnets); err != nil {
		return fmt.Errorf("Error setting subnets: %s", err)
	}

	return nil
}

func getLinkToPublicIpsXML(body io.Reader) (string, error) {
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return "", z.Err()
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if strings.Contains(attr.Val, "PublicIPs") {
							return attr.Val, nil
						}
					}
				}
			}
		}
	}
}
