package overtime

type OvertimeRequest struct {
	Hours float64 `json:"hours" binding:"required,gte=0.5,lte=3"`
}
