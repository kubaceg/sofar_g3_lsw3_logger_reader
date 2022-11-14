package sofar

import "github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"

type Logger struct {
	serialNumber uint
	connPort     ports.CommunicationPort
}

func NewSofarLogger(serialNumber uint, connPort ports.CommunicationPort) *Logger {
	return &Logger{
		serialNumber: serialNumber,
		connPort:     connPort,
	}
}

func (s *Logger) Query() (map[string]interface{}, error) {
	return readData(s.connPort, s.serialNumber)
}

func (s *Logger) Name() string {
	return "sofar"
}
