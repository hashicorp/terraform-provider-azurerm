package validate

import (
	"fmt"
	"strings"
)

var categoryList = []string{
	"abortion",
	"abused-drugs",
	"adult",
	"alcohol-and-tobacco",
	"auctions",
	"business-and-economy",
	"command-and-control",
	"computer-and-internet-info",
	"content-delivery-networks",
	"copyright-infringement",
	"cryptocurrency",
	"dating",
	"dynamic-dns",
	"educational-institutions",
	"entertainment-and-arts",
	"extremism",
	"financial-services",
	"gambling",
	"games",
	"government",
	"grayware",
	"hacking",
	"health-and-medicine",
	"high-risk",
	"home-and-garden",
	"hunting-and-fishing",
	"insufficient-content",
	"internet-communications-and-telephony",
	"internet-portals",
	"job-search",
	"legal",
	"low-risk",
	"malware",
	"medium-risk",
	"military",
	"motor-vehicles",
	"music",
	"newly-registered-domain",
	"news",
	"not-resolved",
	"nudity",
	"online-storage-and-backup",
	"parked",
	"peer-to-peer",
	"personal-sites-and-blogs",
	"philosophy-and-political-advocacy",
	"phishing",
	"private-ip-addresses",
	"proxy-avoidance-and-anonymizers",
	"questionable",
	"real-estate",
	"real-time-detection",
	"recreation-and-hobbies",
	"reference-and-research",
	"religion",
	"search-engines",
	"sex-education",
	"shareware-and-freeware",
	"shopping",
	"social-networking",
	"society",
	"sports",
	"stock-advice-and-tools",
	"streaming-media",
	"swimsuits-and-intimate-apparel",
	"training-and-tools",
	"translation",
	"travel",
	"unknown",
	"weapons",
	"web-advertisements",
	"web-based-email",
	"web-hosting",
}

func CategoryNames(input interface{}, k string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be of type string", k))
		return
	}

	for _, c := range categoryList {
		if strings.EqualFold(v, c) {
			return
		}
	}

	errors = append(errors, fmt.Errorf("%q is not a valid category name", v))

	return
}
