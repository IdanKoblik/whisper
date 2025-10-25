package services

import (
	"errors"
	"testing"
	"whisper-api/communication"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage_Success(t *testing.T) {
	deviceID := "device-123"

	mockConn := &mock.MockConn{}
	communication.Clients = map[string]communication.Conn{
		deviceID: mockConn,
	}

	request := MessageRequest{
		DeviceID:    deviceID,
		Message:     "Hello World",
		Subscribers: []string{"sub1", "sub2"},
	}

	err := SendMessage(request)
	assert.NoError(t, err)
	assert.Equal(t, request, mockConn.LastMessage)
}

func TestSendMessage_DeviceNotConnected(t *testing.T) {
	deviceID := "device-456"
	communication.Clients = map[string]communication.Conn{} // no devices

	request := MessageRequest{
		DeviceID: deviceID,
		Message:  "Hello",
	}

	err := SendMessage(request)
	assert.Error(t, err)
	assert.Equal(t, "device not connected", err.Error())
}

func TestSendMessage_WriteJSONError(t *testing.T) {
	deviceID := "device-789"

	mockConn := &mock.MockConn{ErrToReturn: errors.New("write failed")}
	communication.Clients = map[string]communication.Conn{
		deviceID: mockConn,
	}

	request := MessageRequest{
		DeviceID: deviceID,
		Message:  "Hello",
	}

	err := SendMessage(request)
	assert.Error(t, err)
	assert.Equal(t, "write failed", err.Error())
}
