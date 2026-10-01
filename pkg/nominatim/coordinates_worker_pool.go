package nominatim

type CoordinatesWorkerPool struct {
	client          *Client
	requestsChannel chan CoordinatesRequestToWorker

	pool []*CoordinatesWorker
}

func NewCoordinatesWorkerPool(client *Client) *CoordinatesWorkerPool {
	return &CoordinatesWorkerPool{
		client:          client,
		requestsChannel: make(chan CoordinatesRequestToWorker, 256),
	}
}

// CreateWorker creates a CoordinatesWorker that pulls requests from the
// pool's shared channel. Multiple workers can be created to process the
// same queue concurrently.
func (p *CoordinatesWorkerPool) CreateWorker() *CoordinatesWorker {
	worker := &CoordinatesWorker{
		client:          p.client,
		requestsChannel: p.requestsChannel,
		done:            make(chan struct{}),
	}

	p.pool = append(p.pool, worker)

	return worker
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

func (p *CoordinatesWorkerPool) GetWorkersPool() []*CoordinatesWorker {
	return p.pool
}
