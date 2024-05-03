package xtest

import (
	"net"
	"sync"

	"github.com/stretchr/testify/require"
)

var portLock sync.Mutex

// OpenPort returns an open port.
func OpenPort(t TestingT) int {
	portLock.Lock()
	defer portLock.Unlock()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	require.NoError(t, err)
	return port
}
