package msgraph

import "fmt"

// ListAppAssignments .
func (g *GraphClient) ListAppAssignments(appID string, opts ...ListQueryOption) ([]MobileAppAssignment, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/assignments", appID)
	var marsh struct {
		Target []MobileAppAssignment `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Target, err
}
