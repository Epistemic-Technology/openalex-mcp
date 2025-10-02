package tools

import (
	"context"
	"os"

	"github.com/Sunhill666/goalex"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetWorkQuery struct {
	WorkID string `json:"work_id"`
}

func GetWorkTool() *mcp.Tool {
	inputSchema, err := jsonschema.For[GetWorkQuery](nil)
	if err != nil {
		panic(err)
	}
	getWorkTool := mcp.Tool{
		Name:        "openalex-get-work",
		Description: "Get a work from OpenAlex",
		InputSchema: inputSchema,
	}
	return &getWorkTool
}

func GetWorkToolHandler(ctx context.Context, req *mcp.CallToolRequest, query GetWorkQuery) (*mcp.CallToolResult, *goalex.Work, error) {
	email := os.Getenv("OPENALEX_EMAIL")
	client := goalex.NewClient(goalex.PolitePool(email))
	work, err := client.Works().Get(query.WorkID)
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{}, work, nil
}
