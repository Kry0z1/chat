package chat

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRegistration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		clientConn, serverConn := net.Pipe()
		defer clientConn.Close()
		handler := ServerHandler()
		go handler(serverConn)

		user := UserIn{Username: "testuser"}
		require.NoError(t, json.NewEncoder(clientConn).Encode(user))

		var resp Response
		require.NoError(t, json.NewDecoder(clientConn).Decode(&resp))

		assert.Empty(t, resp.Error)
		assert.Equal(t, user.Username, resp.Username)
	})

	t.Run("bad_json", func(t *testing.T) {
		clientConn, serverConn := net.Pipe()
		defer clientConn.Close()
		handler := ServerHandler()
		go handler(serverConn)

		_, err := clientConn.Write([]byte("{ invalid }"))
		require.NoError(t, err)

		decoder := json.NewDecoder(clientConn)
		var resp Response
		require.NoError(t, decoder.Decode(&resp))
		assert.Equal(t, ErrBadJson.Error(), resp.Error)

		user := UserIn{Username: "testuser"}
		require.NoError(t, json.NewEncoder(clientConn).Encode(user))

		require.NoError(t, decoder.Decode(&resp))
		assert.Equal(t, user.Username, resp.Username)
		assert.Empty(t, resp.Error)
	})
}

func TestTopicCommand(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	handler := ServerHandler()
	go handler(serverConn)

	require.NoError(t, json.NewEncoder(clientConn).Encode(UserIn{Username: "user"}))
	var regResp Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&regResp))

	cmd := Request{Command: "topic", Content: "newtopic"}
	require.NoError(t, json.NewEncoder(clientConn).Encode(cmd))

	pubCmd := Request{Command: "publish", Content: "hello"}
	require.NoError(t, json.NewEncoder(clientConn).Encode(pubCmd))

	var msg Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&msg))

	assert.Equal(t, "newtopic", msg.Topic)
	assert.Equal(t, "hello", msg.Content)

}

func TestCommandClose(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	handler := ServerHandler()
	go handler(serverConn)

	require.NoError(t, json.NewEncoder(clientConn).Encode(UserIn{Username: "user"}))
	var resp Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&resp))

	require.NoError(t, json.NewEncoder(clientConn).Encode(Request{Command: "close"}))

	var closeResp Response
	err := json.NewDecoder(clientConn).Decode(&closeResp)
	assert.Error(t, err)

	_, err = clientConn.Read(make([]byte, 1))
	assert.Error(t, err)
}

func TestCommandUnknown(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	handler := ServerHandler()
	go handler(serverConn)

	require.NoError(t, json.NewEncoder(clientConn).Encode(UserIn{Username: "user"}))
	var regResp Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&regResp))

	cmd := Request{Command: "invalid"}
	require.NoError(t, json.NewEncoder(clientConn).Encode(cmd))

	var resp Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&resp))

	assert.Equal(t, ErrUnknownCommand.Error(), resp.Error)
}

func TestCommandWithBadJSON(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	handler := ServerHandler()
	go handler(serverConn)

	require.NoError(t, json.NewEncoder(clientConn).Encode(UserIn{Username: "user"}))
	var regResp Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&regResp))

	_, err := clientConn.Write([]byte("{ invalid }"))
	require.NoError(t, err)

	var resp Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&resp))

	assert.Equal(t, ErrBadJson.Error(), resp.Error)
}

func TestPublishCommand(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	handler := ServerHandler()
	go handler(serverConn)

	require.NoError(t, json.NewEncoder(clientConn).Encode(UserIn{Username: "user"}))
	var regResp Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&regResp))

	pubCmd := Request{Command: "publish", Content: "test message"}
	require.NoError(t, json.NewEncoder(clientConn).Encode(pubCmd))

	var msg Response
	require.NoError(t, json.NewDecoder(clientConn).Decode(&msg))

	assert.Equal(t, "test message", msg.Content)
	assert.Equal(t, "user", msg.Username)
	assert.Equal(t, "global", msg.Topic)
}
