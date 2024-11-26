package handler

import (
	"net/http"
	usecase "ruantiengo/internal/usecases"
	"time"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	statsService *usecase.StatisticsService
}

func NewStatisticsHandler(statsService *usecase.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statsService: statsService}
}

// GetCompanyStatistics - Rota para estatísticas por empresa
func (h *StatisticsHandler) GetCompanyStatistics(c *gin.Context) {
	tenantID := c.GetHeader("tenantid")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenantid is required in headers"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	branchID := c.Query("branch_id")

	// Validação de datas
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM-DD"})
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use YYYY-MM-DD"})
		return
	}

	if start.After(end) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before end_date"})
		return
	}

	// Chama o serviço para calcular as estatísticas
	stats, err := h.statsService.GetCompanyStatistics(c.Request.Context(), tenantID, branchID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company statistics", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// // GetCustomerStatistics - Rota para estatísticas por cliente
// func (h *StatisticsHandler) GetCustomerStatistics(c *gin.Context) {
// 	tenantID := c.GetHeader("tenant_id")
// 	if tenantID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required in headers"})
// 		return
// 	}

// 	startDate := c.Query("start_date")
// 	endDate := c.Query("end_date")
// 	customerID := c.Query("customer_id")

// 	// Validação de datas
// 	start, err := time.Parse("2006-01-02", startDate)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM-DD"})
// 		return
// 	}
// 	end, err := time.Parse("2006-01-02", endDate)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use YYYY-MM-DD"})
// 		return
// 	}

// 	if start.After(end) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before end_date"})
// 		return
// 	}

// 	// // Chama o serviço para calcular as estatísticas
// 	// stats, err := h.statsService.GetCustomerStatistics(c.Request.Context(), tenantID, customerID, start, end)
// 	// if err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer statistics", "details": err.Error()})
// 	// 	return
// 	// }

// 	c.JSON(http.StatusOK, stats)
//}
