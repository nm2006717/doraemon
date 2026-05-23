package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/frankfoo/doraemon/internal/classify"
	"github.com/frankfoo/doraemon/internal/search"
	"github.com/frankfoo/doraemon/internal/store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewServer(s *store.Store) *server.MCPServer {
	srv := server.NewMCPServer("doraemon", "0.1.0",
		server.WithToolCapabilities(true),
	)

	rec := search.NewRecommender(s)

	srv.AddTool(storeTool(), storeHandler(s))
	srv.AddTool(searchTool(), searchHandler(s))
	srv.AddTool(listTool(), listHandler(s))
	srv.AddTool(recommendTool(), recommendHandler(rec))

	return srv
}

func storeTool() mcp.Tool {
	return mcp.NewTool("store",
		mcp.WithDescription("Store a note, link, or file into the knowledge base"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Title of the entry")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Content to store")),
		mcp.WithString("category", mcp.Description("Category: note, link, file, image. Defaults to note")),
		mcp.WithString("tags", mcp.Description("Comma-separated tags")),
	)
}

func storeHandler(s *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultText("missing required field: title"), nil
		}
		content, err := req.RequireString("content")
		if err != nil {
			return mcp.NewToolResultText("missing required field: content"), nil
		}

		category := req.GetString("category", "")
		tags := req.GetString("tags", "")

		if category == "" || tags == "" {
			result := classify.Classify(title, content)
			if category == "" {
				category = result.Category
			}
			if tags == "" {
				tags = classify.FormatTags(result.Tags)
			}
		}

		id, err := s.Save(&store.Entry{
			Title:    title,
			Content:  content,
			Category: category,
			Tags:     tags,
		})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("failed to store: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Stored successfully with ID: %d (category: %s, tags: %s)", id, category, tags)), nil
	}
}

func searchTool() mcp.Tool {
	return mcp.NewTool("search",
		mcp.WithDescription("Search the knowledge base by keyword"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
	)
}

func searchHandler(s *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := req.RequireString("query")
		if err != nil {
			return mcp.NewToolResultText("missing required field: query"), nil
		}

		ftsQuery := store.BuildFTSQuery(query)
		entries, err := s.Search(ftsQuery)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("search failed: %v", err)), nil
		}

		if len(entries) == 0 {
			return mcp.NewToolResultText("No results found."), nil
		}

		result, _ := json.MarshalIndent(entries, "", "  ")
		return mcp.NewToolResultText(string(result)), nil
	}
}

func listTool() mcp.Tool {
	return mcp.NewTool("list",
		mcp.WithDescription("List all entries or filter by category"),
		mcp.WithString("category", mcp.Description("Filter by category: note, link, file, image. Leave empty for all")),
	)
}

func listHandler(s *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		category := req.GetString("category", "")

		entries, err := s.List(category)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("list failed: %v", err)), nil
		}

		if len(entries) == 0 {
			return mcp.NewToolResultText("No entries found."), nil
		}

		result, _ := json.MarshalIndent(entries, "", "  ")
		return mcp.NewToolResultText(string(result)), nil
	}
}

func recommendTool() mcp.Tool {
	return mcp.NewTool("recommend",
		mcp.WithDescription("Suggest relevant entries based on current conversation context"),
		mcp.WithString("context", mcp.Required(), mcp.Description("Current conversation context or topic to get recommendations for")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of recommendations (default 5)")),
	)
}

func recommendHandler(rec *search.Recommender) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		text, err := req.RequireString("context")
		if err != nil {
			return mcp.NewToolResultText("missing required field: context"), nil
		}

		limit := req.GetInt("limit", 5)

		entries, err := rec.Recommend(text, limit)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("recommend failed: %v", err)), nil
		}

		if len(entries) == 0 {
			return mcp.NewToolResultText("No recommendations found."), nil
		}

		result, _ := json.MarshalIndent(entries, "", "  ")
		return mcp.NewToolResultText(string(result)), nil
	}
}
