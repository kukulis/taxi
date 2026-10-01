package nominatim

import "log"

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
			log.Println("CoordinatesWorker stopped")
			return
		}
	}
}

// Stop signals Handle to return once it has finished processing any
// in-flight request.
func (cw *CoordinatesWorker) Stop() {
	close(cw.done)
}
