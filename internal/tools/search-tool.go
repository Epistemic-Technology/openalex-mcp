package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/Sunhill666/goalex"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchQuery struct {
	Query  string        `json:"query,omitempty"`
	Filter *SearchFilter `json:"filter,omitempty"`
}

type SearchFilter struct {
	AuthorIDs           []string `json:"authorships.author.id,omitempty" jsonschema:"list of OpenAlexIDs"`
	AuthorOrcidIDs      []string `json:"authorships.author.orcid,omitempty" jsonschema:"list of ORCID IDs"`
	MinCitedByCount     int      `json:"min_cited_by_count,omitempty" jsonschema:"minimum number of citations of this work"`
	IsOpenAccess        bool     `json:"open_access.is_oa,omitempty" jsonschema:"whether the work is open access"`
	FromPublicationDate string   `json:"from_publication_date,omitempty" jsonschema:"restrict to publications after this date"`
	ToPublicationDate   string   `json:"to_publication_date,omitempty" jsonschema:"restrict to publications before this date"`
	CitedBy             string   `json:"cited_by,omitempty" jsonschema:"restrict to works cited by this work by OpenAlexID"`
	Cites               string   `json:"cites,omitempty" jsonschema:"restrict to works that cite this work by OpenAlexID"`
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

func SearchToolHandler(ctx context.Context, req *mcp.CallToolRequest, query SearchQuery) (*mcp.CallToolResult, *SearchResult, error) {
	email := os.Getenv("OPENALEX_EMAIL")
	client := goalex.NewClient(goalex.PolitePool(email))

	worksQuery := client.Works().Search(query.Query)

	// Apply filters if provided
	if query.Filter != nil {
		if len(query.Filter.AuthorIDs) > 0 {
			worksQuery = worksQuery.Filter("authorships.author.id", query.Filter.AuthorIDs)
		}
		if len(query.Filter.AuthorOrcidIDs) > 0 {
			worksQuery = worksQuery.Filter("authorships.author.orcid", query.Filter.AuthorOrcidIDs)
		}
		if query.Filter.MinCitedByCount > 0 {
			worksQuery = worksQuery.Filter("cited_by_count", fmt.Sprintf(">=%d", query.Filter.MinCitedByCount))
		}
		if query.Filter.IsOpenAccess {
			worksQuery = worksQuery.Filter("open_access.is_oa", true)
		}
		if query.Filter.FromPublicationDate != "" {
			worksQuery = worksQuery.Filter("from_publication_date", query.Filter.FromPublicationDate)
		}
		if query.Filter.ToPublicationDate != "" {
			worksQuery = worksQuery.Filter("to_publication_date", query.Filter.ToPublicationDate)
		}
		if query.Filter.CitedBy != "" {
			worksQuery = worksQuery.Filter("cited_by", query.Filter.CitedBy)
		}
		if query.Filter.Cites != "" {
			worksQuery = worksQuery.Filter("cites", query.Filter.Cites)
		}
	}

	oaWorks, err := worksQuery.List()
	if err != nil {
		return nil, nil, err
	}

	var works []*CondensedWork
	for _, work := range oaWorks {
		works = append(works, condenseWork(work))
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
