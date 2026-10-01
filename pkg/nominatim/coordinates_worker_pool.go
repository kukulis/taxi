package nominatim

type CoordinatesWorkerPool struct {
	requestsChannel chan CoordinatesRequestToWorker
}

func NewCoordinatesWorkerPool() *CoordinatesWorkerPool {
	return &CoordinatesWorkerPool{
		requestsChannel: make(chan CoordinatesRequestToWorker, 256),
	}
}

// CreateWorker creates a CoordinatesWorker that pulls requests from the
// pool's shared channel. Multiple workers can be created to process the
// same queue concurrently.
func (p *CoordinatesWorkerPool) CreateWorker(client *Client) *CoordinatesWorker {
	return &CoordinatesWorker{
		client:          client,
		requestsChannel: p.requestsChannel,
		done:            make(chan struct{}),
	}
}

// AddRequest enqueues request without blocking. It reports false if the
// queue is full, letting the caller decide how to handle backpressure.
func (p *CoordinatesWorkerPool) AddRequest(request CoordinatesRequestToWorker) bool {
	select {
	case p.requestsChannel <- request:
		return true
	default:
		return false
	}
}
