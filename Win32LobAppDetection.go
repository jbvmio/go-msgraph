package msgraph

type Win32LobAppDetection map[string]interface{}

func (d Win32LobAppDetection) GetODataType() string {
	if v, ok := d[`@odata.type`].(string); ok {
		return v
	}
	return "MISSING"
}

func NewWin32LobAppProductCodeDetection(productCode string) Win32LobAppDetection {
	d := make(Win32LobAppDetection, 3)
	d[`@odata.type`] = `#microsoft.graph.win32LobAppProductCodeDetection`
	d[`productCode`] = productCode
	d[`productVersionOperator`] = `notConfigured`
	return d
}

/*
type Win32LobAppProductCodeDetection struct {
	ODataType              string `json:"@odata.type,omitempty"`
	ProductCode            string `json:"productCode,omitempty"`
	ProductVersionOperator string `json:"productVersionOperator,omitempty"`
	ProductVersion         string `json:"productVersion,omitempty"`
}

func (d *Win32LobAppProductCodeDetection) GetODataType() string {
	return d.ODataType
}

func NewWin32LobAppProductCodeDetection(productCode string) *Win32LobAppProductCodeDetection {
	return &Win32LobAppProductCodeDetection{
		ODataType:              `#microsoft.graph.win32LobAppProductCodeDetection`,
		ProductCode:            productCode,
		ProductVersionOperator: `notConfigured`,
	}
}
*/

// Pwsh:
// https://docs.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapppowershellscriptdetection?view=graph-rest-beta
// Product Code:
// https://docs.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappproductcodedetection?view=graph-rest-beta
