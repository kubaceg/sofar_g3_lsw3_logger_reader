package sofar

import "git.xelasys.ro/sigxcpu/sofar/ports"

type SofarLogger struct {
	serialNumber string
	connPort     ports.CommunicationPort
}

func NewSofarLogger(serialNumber string, connPort ports.CommunicationPort) *SofarLogger {
	return &SofarLogger{
		serialNumber: serialNumber,
		connPort:     connPort,
	}
}

func (s *SofarLogger) Query() (map[string]interface{}, error) {
	reply, err := ReadData(s.connPort, s.serialNumber)

	return reply, err
}

func (s *SofarLogger) Name() string {
	return "sofar"
}
