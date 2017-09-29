package birdsocket

import (
	"net"
	"regexp"
	"strings"
)

var birdReturnCodeRegex *regexp.Regexp

func init() {
	birdReturnCodeRegex = regexp.MustCompile("\\d{4} \n$")
}

// BirdSocket encapsulates communication with Bird routing daemon
type BirdSocket struct {
	SocketPath string
	BufferSize int
	Timeout    int
	conn       net.Conn
}

// NewSocket creates a new socket
func NewSocket(socketPath string) *BirdSocket {
	return &BirdSocket{SocketPath: socketPath, BufferSize: 4096, Timeout: 30}
}

// Query sends an ad hoc query to Bird and waits for the reply
func Query(socketPath, qry string) ([]byte, error) {
	s := NewSocket(socketPath)
	_, err := s.Connect()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	return s.Query(qry)
}

// Connect connects to the Bird unix socket
func (s *BirdSocket) Connect() ([]byte, error) {
	var err error
	s.conn, err = net.Dial("unix", s.SocketPath)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, s.BufferSize)
	n, err := s.conn.Read(buf[:])
	if err != nil {
		return nil, err
	}

	return buf[:n], err
}

// Close closes the connection to the socket
func (s *BirdSocket) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

// Query sends an query to Bird and waits for the reply
func (s *BirdSocket) Query(qry string) ([]byte, error) {
	_, err := s.conn.Write([]byte(strings.Trim(qry, "\n") + "\n"))
	if err != nil {
		return nil, err
	}

	output, err := s.readFromSocket(s.conn)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *BirdSocket) readFromSocket(conn net.Conn) ([]byte, error) {
	b := make([]byte, 0)
	buf := make([]byte, s.BufferSize)

	done := false
	for !done {
		n, err := conn.Read(buf[:])
		if err != nil {
			return nil, err
		}

		b = append(b, buf[:n]...)
		done = endsWithBirdReturnCode(b)
	}

	return b, nil
}

func endsWithBirdReturnCode(b []byte) bool {
	if len(b) < 6 {
		return false
	}

	return birdReturnCodeRegex.Match(b[len(b)-6:])
}
