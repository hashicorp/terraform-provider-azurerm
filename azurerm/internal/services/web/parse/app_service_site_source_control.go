package parse

// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-group/providers/Microsoft.Web/sites/test-app/sourcecontrols/web

type AppServiceSiteSourceControlId struct {
	ResourceGroup string
	AppName       string
	Name          string
}

func AppServiceSiteSourceControlID(input string) (*AppServiceSiteSourceControlId, error) {
	return &AppServiceSiteSourceControlId{}, nil
}
