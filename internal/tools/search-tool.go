package tools

import (
	"context"
	"os"

	"github.com/Sunhill666/goalex"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchQuery struct {
	Query string `json:"query"`
}

func SearchTool() *mcp.Tool {
	inputschema, err := jsonschema.For[SearchQuery](nil)
	if err != nil {
		panic(err)
	}
	searchTool := mcp.Tool{
		Name:        "openalex-search",
		Description: "Search OpenAlex for works",
		InputSchema: inputschema,
	}
	return &searchTool
}

type SearchResult struct {
	Works []*goalex.Work `json:"works"`
}

func SearchToolHandler(ctx context.Context, req *mcp.CallToolRequest, query SearchQuery) (*mcp.CallToolResult, *SearchResult, error) {
	email := os.Getenv("OPENALEX_EMAIL")
	client := goalex.NewClient(goalex.PolitePool(email))
	works, err := client.Works().Search(query.Query).List()
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{}, &SearchResult{Works: works}, nil
}
