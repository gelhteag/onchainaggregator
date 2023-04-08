package uniswap

import (
	"testing"
)

func TestBuildArgs(t *testing.T) {
	args := map[string]interface{}{
		"id":   "0x6B175474E89094C44Da98b954EedeAC495271d0F",
		"skip": 100,
		"sort": map[string]interface{}{
			"field": "name",
			"order": "ASC",
		},
	}

	argsString := BuildArgs(args)
	expectedArgsString := `id: "0x6B175474E89094C44Da98b954EedeAC495271d0F", skip: 100, sort: {field: "name", order: "ASC"}`
	if argsString != expectedArgsString {
		t.Errorf("Unexpected args string:\nGot:      %s\nExpected: %s", argsString, expectedArgsString)
	}
}

func TestBuildFields(t *testing.T) {
	type Token struct {
		ID     string `graphql:"id"`
		Name   string `graphql:"name"`
		Symbol string `graphql:"symbol"`
	}

	fieldsString := BuildFields(&Token{})
	expectedFieldsString := "id\nname\nsymbol"
	if fieldsString != expectedFieldsString {
		t.Errorf("Unexpected fields string: %s", fieldsString)
	}
}
func TestRunGraphQLQuery(t *testing.T) {
	query := `
		query {
		  token(id: "0x6b175474e89094c44da98b954eedeac495271d0f") {
		    id
		    symbol
		    name
		  }
		}
	`
	response, err := RunGraphQLQuery(query)
	if err != nil {
		t.Fatalf("Error executing GraphQL query: %v", err)
	}

	tokenData := response["token"].(map[string]interface{})
	expectedID := "0x6b175474e89094c44da98b954eedeac495271d0f"
	expectedSymbol := "DAI"
	expectedName := "Dai Stablecoin"

	if tokenData["id"] != expectedID {
		t.Errorf("Expected ID %s but got %v", expectedID, tokenData["id"])
	}

	if tokenData["symbol"] != expectedSymbol {
		t.Errorf("Expected symbol %s but got %v", expectedSymbol, tokenData["symbol"])
	}

	if tokenData["name"] != expectedName {
		t.Errorf("Expected name %s but got %v", expectedName, tokenData["name"])
	}
}
