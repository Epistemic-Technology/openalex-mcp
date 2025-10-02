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
	Works []*CondensedWork `json:"works"`
}

type CondensedWork struct {
	Abstract        string             `json:"abstract,omitempty"`
	Authorships     []*CondensedAuthor `json:"authorships,omitempty"`
	CitedByCount    int                `json:"cited_by_count,omitempty"`
	CreatedDate     string             `json:"created_date,omitempty"`
	DOI             string             `json:"doi,omitempty"`
	ID              string             `json:"id,omitempty"`
	PublicationDate string             `json:"publication_date,omitempty"`
	PublicationYear int                `json:"publication_year,omitempty"`
	Title           string             `json:"title,omitempty"`
	Type            string             `json:"type,omitempty"`
}

type CondensedAuthor struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	ORCID       string `json:"orcid,omitempty"`
}

func SearchToolHandler(ctx context.Context, req *mcp.CallToolRequest, query SearchQuery) (*mcp.CallToolResult, *SearchResult, error) {
	email := os.Getenv("OPENALEX_EMAIL")
	client := goalex.NewClient(goalex.PolitePool(email))
	oaWorks, err := client.Works().Search(query.Query).List()
	var works []*CondensedWork
	for _, work := range oaWorks {
		works = append(works, condenseWork(work))
	}
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{}, &SearchResult{Works: works}, nil
}

func condenseWork(work *goalex.Work) *CondensedWork {
	var condensedAuthorships []*CondensedAuthor
	for _, authorship := range work.Authorships {
		condensedAuthorships = append(condensedAuthorships, &CondensedAuthor{
			ID:          authorship.Author.ID,
			DisplayName: authorship.Author.DisplayName,
			ORCID:       authorship.Author.ORCID,
		})
	}
	return &CondensedWork{
		Abstract:        work.Abstract,
		Authorships:     condensedAuthorships,
		CitedByCount:    work.CitedByCount,
		CreatedDate:     work.CreatedDate,
		DOI:             work.DOI,
		ID:              work.ID,
		PublicationDate: work.PublicationDate,
		PublicationYear: work.PublicationYear,
		Title:           work.Title,
		Type:            work.Type,
	}
}
