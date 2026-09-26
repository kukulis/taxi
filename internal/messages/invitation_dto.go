package messages

type InvitationDto struct {
	InvitationId string `json:"invitation_id"`
	// Status invitation status
	Status   string `json:"status"`
	DriverId string `json:"driver_id"`
	Lat      string `json:"lat"`
	Lon      string `json:"lon"`
}
