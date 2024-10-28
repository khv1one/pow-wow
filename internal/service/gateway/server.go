package gateway

type Challenger interface {
	MakeChallenge() (string, error)
}

type Supervisor interface {
	Oversee(challenge, nonce string) (string, error)
}

type Server struct {
	challenger Challenger
	supervisor Supervisor
}

func NewServer(c Challenger, s Supervisor) *Server {
	return &Server{challenger: c, supervisor: s}
}
