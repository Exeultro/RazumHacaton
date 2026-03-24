package handlers

import (
	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	searchType := c.DefaultQuery("type", "all")

	if query == "" {
		utils.BadRequest(c, "Search query is required")
		return
	}

	results, err := h.searchService.Search(query, searchType)
	if err != nil {
		utils.InternalServerError(c, "Failed to search: "+err.Error())
		return
	}

	utils.SuccessResponse(c, results)
}
