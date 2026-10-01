package nominatim

type AddressWorker struct {
	client          *Client
	requestsChannel chan AddressRequestToWorker
	done            chan struct{}
}

func (aw *AddressWorker) Handle() {
	for {
		select {
		case request := <-aw.requestsChannel:
			address, err := aw.client.ReverseGeocode(request.Lat, request.Lon)
			request.Callback(address, err)
		case <-aw.done:
			return
		}
	}
}

// Stop signals Handle to return once it has finished processing any
// in-flight request. Requests still queued in requestsChannel are dropped.
func (aw *AddressWorker) Stop() {
	close(aw.done)
}
