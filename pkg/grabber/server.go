// Package grabber implements start and stop grabbers.
package grabber

type Grabber interface {
	Start(chan bool)
}

// Server -.
type Server struct {
	shutdown    chan bool
	apiGrabbers []Grabber
}

// New -.
func New(apiGrabbers []Grabber) *Server {
	s := &Server{shutdown: make(chan bool, 1), apiGrabbers: apiGrabbers}
	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		for _, grabber := range s.apiGrabbers {
			grabber.Start(s.shutdown)
		}
	}()
}

// Shutdown -.
func (s *Server) Shutdown() {
	s.shutdown <- true
}
