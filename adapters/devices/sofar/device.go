package sofar

import "github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"

type Logger struct {
	serialNumber  uint
	connPort      ports.CommunicationPort
	attrWhiteList map[string]struct{}
	attrBlackList []string
}

// for a set in go we use a map of keys -> empty struct
func toSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{}, len(slice))
	v := struct{}{}
	for _, s := range slice {
		set[s] = v
	}
	return set
}

func NewSofarLogger(serialNumber uint, connPort ports.CommunicationPort, attrWhiteList []string, attrBlackList []string) *Logger {
	return &Logger{
		serialNumber:  serialNumber,
		connPort:      connPort,
		attrWhiteList: toSet(attrWhiteList),
		attrBlackList: attrBlackList,
	}
}

func (s *Logger) Query() (map[string]interface{}, error) {
	return readData(s.connPort, s.serialNumber, s.attrWhiteList, s.attrBlackList)
}

func (s *Logger) Name() string {
	return "sofar"
}
