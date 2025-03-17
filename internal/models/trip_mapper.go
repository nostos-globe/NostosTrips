package models

type TripMapper struct{}

func (m *TripMapper) ToTrip(req struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}, tokenResponse interface{}) Trip {
	return Trip{
		Name:        req.Name,
		Description: req.Description,
		Visibility:  req.Visibility,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		UserID:      tokenResponse.(uint),
	}
}
func (m *TripMapper) ToTripUpdate(req struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}, tokenResponse interface{}) Trip {
	return Trip{
		TripID:      req.ID,
		Name:        req.Name,
		Description: req.Description,
		Visibility:  req.Visibility,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		UserID:      tokenResponse.(uint),
	}
}
