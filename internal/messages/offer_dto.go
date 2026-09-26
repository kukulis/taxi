package messages

type OfferDto struct {
	InvitationId string `json:"invitation_id"`
	// Status invitation status
	Status      string `json:"status"`
	PassengerId string `json:"passenger_id"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
}

func GetOfferPassengerId(dto *OfferDto) string { return dto.PassengerId }
