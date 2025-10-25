package mock

type MockConn struct {
	LastMessage interface{}
	ErrToReturn error
}

func (m *MockConn) WriteJSON(v interface{}) error {
	m.LastMessage = v
	return m.ErrToReturn
}
