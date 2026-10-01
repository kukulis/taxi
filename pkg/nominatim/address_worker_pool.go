package nominatim

type AddressWorkerPool struct {
	client          *Client
	requestsChannel chan AddressRequestToWorker
	pool            []*AddressWorker
}

func NewAddressWorkerPool(client *Client) *AddressWorkerPool {
	return &AddressWorkerPool{
		client:          client,
		requestsChannel: make(chan AddressRequestToWorker, 256),
	}
}

// CreateWorker creates an AddressWorker that pulls requests from the pool's
// shared channel. Multiple workers can be created to process the same queue
// concurrently.
func (p *AddressWorkerPool) CreateWorker() *AddressWorker {
	worker := &AddressWorker{
		client:          p.client,
		requestsChannel: p.requestsChannel,
		done:            make(chan struct{}),
	}

	p.pool = append(p.pool, worker)

	return worker
}

// AddRequest enqueues request without blocking. It reports false if the
// queue is full, letting the caller decide how to handle backpressure.
func (p *AddressWorkerPool) AddRequest(request AddressRequestToWorker) bool {
	select {
	case p.requestsChannel <- request:
		return true
	default:
		return false
	}
}

func (p *AddressWorkerPool) GetWorkersPool() []*AddressWorker {
	return p.pool
}
