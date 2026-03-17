package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Input types ---

type ListPortalsInput struct{}

type UpdatePortalInput struct {
	PortalID int    `json:"portal_id"`
	Name     string `json:"name,omitempty"`
	Slug     string `json:"slug,omitempty"`
}

type ListArticlesInput struct {
	PortalID int `json:"portal_id"`
}

type CreateArticleInput struct {
	PortalID    int    `json:"portal_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	CategoryID  *int   `json:"category_id,omitempty"`
	AuthorID    *int   `json:"author_id,omitempty"`
}

type UpdateArticleInput struct {
	PortalID    int    `json:"portal_id"`
	ArticleID   int    `json:"article_id"`
	Title       string `json:"title,omitempty"`
	Content     string `json:"content,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	CategoryID  *int   `json:"category_id,omitempty"`
}

type DeleteArticleInput struct {
	PortalID  int `json:"portal_id"`
	ArticleID int `json:"article_id"`
}

type ListCategoriesInput struct {
	PortalID int `json:"portal_id"`
}

type CreateCategoryInput struct {
	PortalID    int    `json:"portal_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Locale      string `json:"locale,omitempty"`
	Position    *int   `json:"position,omitempty"`
	ParentID    *int   `json:"parent_id,omitempty"`
}

type UpdateCategoryInput struct {
	PortalID    int    `json:"portal_id"`
	CategoryID  int    `json:"category_id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Locale      string `json:"locale,omitempty"`
	Position    *int   `json:"position,omitempty"`
}

type DeleteCategoryInput struct {
	PortalID   int `json:"portal_id"`
	CategoryID int `json:"category_id"`
}

// RegisterHelpCenterTools registers help center (portals, articles, categories) tools.
func RegisterHelpCenterTools(server *mcp.Server, client *chatwoot.Client) {

	// --- list_portals ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_portals",
		Description: "List all help center portals in the account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListPortalsInput) (*mcp.CallToolResult, any, error) {
		portals, err := client.ListPortals(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, p := range portals {
			sb.WriteString(fmt.Sprintf("- [%d] %s (slug: %s)\n", p.ID, p.Name, p.Slug))
		}
		if sb.Len() == 0 {
			sb.WriteString("No portals found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- update_portal ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_portal",
		Description: "Update a help center portal's name or slug.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdatePortalInput) (*mcp.CallToolResult, any, error) {
		updateReq := chatwoot.UpdatePortalRequest{}
		if input.Name != "" {
			updateReq.Name = &input.Name
		}
		if input.Slug != "" {
			updateReq.Slug = &input.Slug
		}
		portal, err := client.UpdatePortal(ctx, input.PortalID, updateReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Portal #%d updated! Name: %s, Slug: %s", portal.ID, portal.Name, portal.Slug)), nil, nil
	})

	// --- list_articles ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_articles",
		Description: "List all articles in a help center portal.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListArticlesInput) (*mcp.CallToolResult, any, error) {
		articles, err := client.ListArticles(ctx, input.PortalID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Articles in portal #%d:\n\n", input.PortalID))
		for _, a := range articles {
			sb.WriteString(fmt.Sprintf("- [%d] %s (status: %s)\n", a.ID, a.Title, a.Status))
		}
		if len(articles) == 0 {
			sb.WriteString("No articles found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- create_article ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_article",
		Description: "Create a new article in a help center portal. Requires portal_id, title, and content.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateArticleInput) (*mcp.CallToolResult, any, error) {
		article, err := client.CreateArticle(ctx, input.PortalID, chatwoot.CreateArticleRequest{
			Title:       input.Title,
			Content:     input.Content,
			Description: input.Description,
			Status:      input.Status,
			CategoryID:  input.CategoryID,
			AuthorID:    input.AuthorID,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Article created! ID: %d, Title: %s", article.ID, article.Title)), nil, nil
	})

	// --- update_article ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_article",
		Description: "Update an article in a help center portal. Provide only fields you want to change.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateArticleInput) (*mcp.CallToolResult, any, error) {
		updateReq := chatwoot.UpdateArticleRequest{
			CategoryID: input.CategoryID,
		}
		if input.Title != "" {
			updateReq.Title = &input.Title
		}
		if input.Content != "" {
			updateReq.Content = &input.Content
		}
		if input.Description != "" {
			updateReq.Description = &input.Description
		}
		if input.Status != "" {
			updateReq.Status = &input.Status
		}
		article, err := client.UpdateArticle(ctx, input.PortalID, input.ArticleID, updateReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Article #%d updated! Title: %s", article.ID, article.Title)), nil, nil
	})

	// --- delete_article ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_article",
		Description: "Delete an article from a help center portal.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteArticleInput) (*mcp.CallToolResult, any, error) {
		if err := client.DeleteArticle(ctx, input.PortalID, input.ArticleID); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Article #%d deleted from portal #%d.", input.ArticleID, input.PortalID)), nil, nil
	})

	// --- list_categories ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_categories",
		Description: "List all categories in a help center portal.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListCategoriesInput) (*mcp.CallToolResult, any, error) {
		categories, err := client.ListCategories(ctx, input.PortalID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Categories in portal #%d:\n\n", input.PortalID))
		for _, c := range categories {
			sb.WriteString(fmt.Sprintf("- [%d] %s (slug: %s)\n", c.ID, c.Name, c.Slug))
			if c.Description != "" {
				sb.WriteString(fmt.Sprintf("    %s\n", c.Description))
			}
		}
		if len(categories) == 0 {
			sb.WriteString("No categories found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- create_category ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_category",
		Description: "Create a new category in a help center portal. Requires portal_id and name.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateCategoryInput) (*mcp.CallToolResult, any, error) {
		category, err := client.CreateCategory(ctx, input.PortalID, chatwoot.CreateCategoryRequest{
			Name:        input.Name,
			Slug:        input.Slug,
			Description: input.Description,
			Locale:      input.Locale,
			Position:    input.Position,
			ParentID:    input.ParentID,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Category created! ID: %d, Name: %s", category.ID, category.Name)), nil, nil
	})

	// --- update_category ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_category",
		Description: "Update a category in a help center portal. Provide only fields you want to change.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateCategoryInput) (*mcp.CallToolResult, any, error) {
		updateReq := chatwoot.UpdateCategoryRequest{
			Position: input.Position,
		}
		if input.Name != "" {
			updateReq.Name = &input.Name
		}
		if input.Description != "" {
			updateReq.Description = &input.Description
		}
		if input.Locale != "" {
			updateReq.Locale = &input.Locale
		}
		category, err := client.UpdateCategory(ctx, input.PortalID, input.CategoryID, updateReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Category #%d updated! Name: %s", category.ID, category.Name)), nil, nil
	})

	// --- delete_category ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_category",
		Description: "Delete a category from a help center portal.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteCategoryInput) (*mcp.CallToolResult, any, error) {
		if err := client.DeleteCategory(ctx, input.PortalID, input.CategoryID); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Category #%d deleted from portal #%d.", input.CategoryID, input.PortalID)), nil, nil
	})
}
