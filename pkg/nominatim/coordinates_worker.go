package nominatim

type CoordinatesWorker struct {
	client          *Client
	requestsChannel chan CoordinatesRequestToWorker
	done            chan struct{}
}

func (cw *CoordinatesWorker) Handle() {
	for {
		select {
		case request := <-cw.requestsChannel:
			lat, lon, err := cw.client.Geocode(request.Address)
			request.Callback(lat, lon, err)
		case <-cw.done:
			return
		}
	}
}

// Stop signals Handle to return once it has finished processing any
// in-flight request. Requests still queued in requestsChannel are dropped.
func (cw *CoordinatesWorker) Stop() {
	close(cw.done)
}
