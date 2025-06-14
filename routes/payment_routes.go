// routes/payment_routes.go
package routes

import (
	"github.com/adipras/tirta-saas-backend/controllers"
	"github.com/adipras/tirta-saas-backend/middleware"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(r *gin.Engine) {
	group := r.Group("/api/payments")
	group.Use(middleware.JWTAuthMiddleware(), middleware.AdminOnly())

	group.POST("/", controllers.CreatePayment)
	group.GET("/customer/:customer_id", controllers.GetPaymentHistoryByCustomerID)

}
