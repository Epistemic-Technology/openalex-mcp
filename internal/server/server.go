package server

import (
	"github.com/Epistemic-Technology/openalex-mcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func CreateServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "openalex-mcp", Version: "v0.0.1"}, nil)
	mcp.AddTool(server, tools.SearchTool(), tools.SearchToolHandler)
	mcp.AddTool(server, tools.GetWorkTool(), tools.GetWorkToolHandler)
	return server
}
