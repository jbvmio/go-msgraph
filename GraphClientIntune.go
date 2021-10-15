package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ListWin32LobApps .
func (g *GraphClient) ListWin32LobApps(opts ...ListQueryOption) ([]Win32LobApp, error) {
	resource := "/deviceAppManagement/mobileApps"
	var marsh struct {
		Apps []Win32LobApp `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Apps, err
}

// ListWin32LobAppContentVersions .
func (g *GraphClient) ListWin32LobAppContentVersions(identifier string, opts ...ListQueryOption) ([]TypeAndID, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions", identifier)
	var marsh struct {
		Versions []TypeAndID `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Versions, err
}

// ListWin32LobAppContentFiles lists properties and relationships of the mobileAppContentFile objects.
func (g *GraphClient) ListWin32LobAppContentFiles(identifier, version string, opts ...ListQueryOption) ([]MobileAppContentFile, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files", identifier, version)
	var marsh struct {
		Files []MobileAppContentFile `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Files, err
}

// GetWin32LobApp .
func (g *GraphClient) GetWin32LobApp(identifier string, opts ...GetQueryOption) (Win32LobApp, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s", identifier)
	app := Win32LobApp{}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &app)
	return app, err
}

// CreateWin32LobApp .
func (g *GraphClient) CreateWin32LobApp(app Win32LobApp, opts ...CreateQueryOption) (Win32LobApp, error) {
	resource := "/deviceAppManagement/mobileApps"
	bodyBytes, err := json.Marshal(app)
	if err != nil {
		return app, err
	}
	reader := bytes.NewReader(bodyBytes)
	err = g.makePOSTAPICall(resource, compileCreateQueryOptions(opts), reader, &app)
	return app, err
}

// CreateWin32LobAppContentVersion creates a new Content Version for the given Win32LobApp ID and returns
// the Newly created ContentVersion ID string and any errors.
func (g *GraphClient) CreateWin32LobAppContentVersion(identifier string, opts ...CreateQueryOption) (TypeAndID, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions", identifier)
	var ID TypeAndID
	err := g.makePOSTAPICall(resource, compileCreateQueryOptions(opts), bytes.NewReader([]byte(`{}`)), &ID)
	return ID, err
}
