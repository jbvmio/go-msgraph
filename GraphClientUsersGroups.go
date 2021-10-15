package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ListUsers returns a list of all users
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user_list
func (g *GraphClient) ListUsers(opts ...ListQueryOption) (Users, error) {
	resource := "/users"
	var marsh struct {
		Users Users `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	marsh.Users.setGraphClient(g)
	return marsh.Users, err
}

// ListGroups returns a list of all groups
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/group_list
func (g *GraphClient) ListGroups(opts ...ListQueryOption) (Groups, error) {
	resource := "/groups"

	var reqParams = compileListQueryOptions(opts)

	var marsh struct {
		Groups Groups `json:"value"`
	}
	err := g.makeGETAPICall(resource, reqParams, &marsh)
	marsh.Groups.setGraphClient(g)
	return marsh.Groups, err
}

// GetUser returns the user object associated to the given user identified by either
// the given ID or userPrincipalName
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user_get
func (g *GraphClient) GetUser(identifier string, opts ...GetQueryOption) (User, error) {
	resource := fmt.Sprintf("/users/%v", identifier)
	user := User{graphClient: g}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &user)
	return user, err
}

// GetGroup returns the group object identified by the given groupID.
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/group_get
func (g *GraphClient) GetGroup(groupID string, opts ...GetQueryOption) (Group, error) {
	resource := fmt.Sprintf("/groups/%v", groupID)
	group := Group{graphClient: g}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &group)
	return group, err
}

// CreateUser creates a new user given a user object and returns and updated object
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user-post-users
func (g *GraphClient) CreateUser(userInput User, opts ...CreateQueryOption) (User, error) {
	user := User{graphClient: g}
	bodyBytes, err := json.Marshal(userInput)
	if err != nil {
		return user, err
	}

	reader := bytes.NewReader(bodyBytes)
	err = g.makePOSTAPICall("/users", compileCreateQueryOptions(opts), reader, &user)

	return user, err
}
