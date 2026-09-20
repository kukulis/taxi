package handler

import (
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
)

type TaxiMessageHandler struct {
	mainState *state.MainState
}

func NewTaxiMessageHandler(mainState *state.MainState) *TaxiMessageHandler {
	return &TaxiMessageHandler{mainState: mainState}
}

func (handler *TaxiMessageHandler) Handle(messageChannel <-chan ws.ClientMessage) {
	// TODO
}
