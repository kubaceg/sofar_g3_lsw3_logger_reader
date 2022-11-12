package sofar

import "github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"

type SofarLogger struct {
	serialNumber uint
	connPort     ports.CommunicationPort
}

func NewSofarLogger(serialNumber uint, connPort ports.CommunicationPort) *SofarLogger {
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
