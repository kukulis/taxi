package nominatim

type AddressResponseCallback func(address string, err error)

type AddressRequestToWorker struct {
	Lat      float64
	Lon      float64
	Callback AddressResponseCallback
}

func NewAddressRequestToWorker(lat float64, lon float64, callback AddressResponseCallback) AddressRequestToWorker {
	return AddressRequestToWorker{
		Lat:      lat,
		Lon:      lon,
		Callback: callback,
	}
}
