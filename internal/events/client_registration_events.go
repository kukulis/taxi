package events

const (
	ClientRegisteredEventName   = "ClientRegistered"
	ClientUnregisteredEventName = "ClientUnregistered"
	ClientTypeDriver            = "driver"
	ClientTypePassenger         = "passenger"
)

type ClientRegisteredEvent struct {
	ClientId   string
	ClientType string
}

type ClientUnregisteredEvent struct {
	ClientId   string
	ClientType string
}

// ClientRegisteredOption sets an optional field on a ClientRegisteredEvent.
type ClientRegisteredOption func(*ClientRegisteredEvent)

func WithRegisteredClientType(clientType string) ClientRegisteredOption {
	return func(e *ClientRegisteredEvent) { e.ClientType = clientType }
}

func NewClientRegisteredEvent(clientId string, opts ...ClientRegisteredOption) *ClientRegisteredEvent {
	e := &ClientRegisteredEvent{
		ClientId: clientId,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// ClientUnregisteredOption sets an optional field on a ClientUnregisteredEvent.
type ClientUnregisteredOption func(*ClientUnregisteredEvent)

func WithUnregisteredClientType(clientType string) ClientUnregisteredOption {
	return func(e *ClientUnregisteredEvent) { e.ClientType = clientType }
}

func NewClientUnregisteredEvent(clientId string, opts ...ClientUnregisteredOption) *ClientUnregisteredEvent {
	e := &ClientUnregisteredEvent{
		ClientId: clientId,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *ClientRegisteredEvent) GetName() string {
	return ClientRegisteredEventName
}

func (e *ClientUnregisteredEvent) GetName() string {
	return ClientUnregisteredEventName
}
