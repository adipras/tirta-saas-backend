package requests

type CreateSubscriptionTypeRequest struct {
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	RegistrationFee float64 `json:"registration_fee" binding:"required"`
	MonthlyFee      float64 `json:"monthly_fee" binding:"required"`
	MaintenanceFee  float64 `json:"maintenance_fee"`
	LateFeePerDay   float64 `json:"late_fee_per_day"`
	MaxLateFee      float64 `json:"max_late_fee"`
}
