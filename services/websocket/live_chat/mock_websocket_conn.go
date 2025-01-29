package live_chat

import "encoding/json"

// MockWebSocketConn is a mock implementation of the WebSocketConn interface.
type MockWebSocketConn struct {
	ReadJSONData string
	ReadJSONErr  error
	WriteJSONErr error
	LastMessage  interface{}
}

func (m *MockWebSocketConn) ReadJSON(v interface{}) error {
	if m.ReadJSONErr != nil {
		return m.ReadJSONErr
	}
	return json.Unmarshal([]byte(m.ReadJSONData), v)
}

func (m *MockWebSocketConn) WriteJSON(v interface{}) error {
	if m.WriteJSONErr != nil {
		return m.WriteJSONErr
	}
	m.LastMessage = v
	return nil
}

func (m *MockWebSocketConn) Close() error {
	return nil
}
