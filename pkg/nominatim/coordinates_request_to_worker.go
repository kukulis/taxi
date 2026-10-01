package nominatim

type CoordinatesResponseCallback func(Lat float64, Lon float64, err error)

type CoordinatesRequestToWorker struct {
	Address  string
	Callback CoordinatesResponseCallback
}

func NewCoordinatesRequestToWorker(address string, callback CoordinatesResponseCallback) *CoordinatesRequestToWorker {
	return &CoordinatesRequestToWorker{
		Address:  address,
		Callback: callback,
	}
}
