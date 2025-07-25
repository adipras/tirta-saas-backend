package routes

import (
	"github.com/adipras/tirta-saas-backend/controllers"
	"github.com/adipras/tirta-saas-backend/middleware"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(r *gin.Engine) {
	group := r.Group("/api/customers")
	group.Use(middleware.JWTAuthMiddleware(), middleware.AdminOnly())

	r.POST("/invoices/generate-monthly", controllers.GenerateMonthlyInvoice)
}
