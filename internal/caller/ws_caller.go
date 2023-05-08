package caller

type wsCaller struct {
	host string
}

// websocket call
func (w *wsCaller) call(method string, params, reply any) error {
	return nil
}

func (w *wsCaller) close() error {
	return nil
}
