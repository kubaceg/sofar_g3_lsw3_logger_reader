package sofar

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"git.xelasys.ro/sigxcpu/sofar/ports"
	"github.com/sigurn/crc16"
)

type LSWRequest struct {
	serialNumber  string
	startRegister int
	endRegister   int
	frameBytes    []byte
}

func NewLSWRequest(serialNumber string, startRegister int, endRegister int) LSWRequest {
	return LSWRequest{
		serialNumber:  serialNumber,
		startRegister: startRegister,
		endRegister:   endRegister,
		frameBytes:    make([]byte, 36),
	}
}

func (l LSWRequest) ToBytes() []byte {
	buf := l.frameBytes

	// preamble
	buf[0] = 0xa5
	binary.BigEndian.PutUint16(buf[1:], 0x1700)
	binary.BigEndian.PutUint16(buf[3:], 0x1045)
	buf[5] = 0x00
	buf[6] = 0x00

	// convert serial number to unsigned 32-bit int

	uint32SerialNumber, _ := strconv.Atoi(l.serialNumber)
	// fmt.Printf("serial number: %0X\n", uint32SerialNumber)
	binary.LittleEndian.PutUint32(buf[7:], uint32(uint32SerialNumber))

	buf[11] = 0x02

	binary.BigEndian.PutUint16(buf[26:], 0x0103)

	binary.BigEndian.PutUint16(buf[28:], uint16(l.startRegister))
	binary.BigEndian.PutUint16(buf[30:], uint16(l.endRegister-l.startRegister+1))

	// compute crc
	table := crc16.MakeTable(crc16.CRC16_MODBUS)
	modbusCRC := crc16.Checksum(buf[26:32], table)

	// h := crc16.New(crc16.Modbus)
	// h.Reset()
	// h.Write(buf[26:32])
	// in := make([]byte, 0, h.Size())
	// in = h.Sum(in)
	// modbusCRC := binary.BigEndian.Uint16(in[0:2])

	// append crc
	binary.LittleEndian.PutUint16(buf[32:], modbusCRC)

	// compute & append frame crc
	buf[34] = l.Checksum()
	// buf.Write([]byte{l.Checksum()})

	// end of frame
	buf[35] = 0x15

	return buf

}

func (l LSWRequest) String() string {
	return fmt.Sprintf("% 0X", l.ToBytes())
}

func (l LSWRequest) Checksum() uint8 {
	var checksum uint8
	for _, b := range l.frameBytes[1 : len(l.frameBytes)-2] {
		checksum += b
	}
	return checksum
}

func ReadData(connPort ports.CommunicationPort, serialNumber string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	reply, err := readData(rrGridOutput, connPort, serialNumber)
	if err != nil {
		return nil, err
	}

	for k, v := range reply {
		result[k] = v
	}

	reply, err = readData(rrPVOutput, connPort, serialNumber)
	if err != nil {
		return nil, err
	}

	for k, v := range reply {
		result[k] = v
	}

	reply, err = readData(rrPVGeneration, connPort, serialNumber)
	if err != nil {
		return nil, err
	}

	for k, v := range reply {
		result[k] = v
	}

	reply, err = readData(rrSystemInfo, connPort, serialNumber)
	if err != nil {
		return nil, err
	}

	for k, v := range reply {
		result[k] = v
	}

	return result, err
}

func readData(rr RegisterRange, connPort ports.CommunicationPort, serialNumber string) (map[string]interface{}, error) {

	lswRequest := NewLSWRequest(serialNumber, rr.start, rr.end)

	commandBytes := lswRequest.ToBytes()

	err := connPort.Open()
	if err != nil {
		return nil, err
	}

	defer connPort.Close()

	// send the command
	_, err = connPort.Write(commandBytes)
	if err != nil {
		return nil, err
	}

	// read the result
	buf := make([]byte, 2048)
	n, err := connPort.Read(buf)
	if err != nil {
		return nil, err
	}
	// truncate the buffer
	buf = buf[:n]
	if len(buf) < 48 {
		// short reply
		return nil, fmt.Errorf("short reply: %d bytes", n)
	}

	replyBytesCount := buf[27]

	modbusReply := buf[28 : 28+replyBytesCount]

	// shove the data into the reply
	reply := make(map[string]interface{})

	for _, f := range rr.replyFields {
		fieldOffset := (f.register - rr.start) * 2

		if fieldOffset > len(modbusReply)-2 {
			// skip invalid offset
			continue
		}

		switch f.valueType {
		case "U16":
			reply[f.name] = binary.BigEndian.Uint16(modbusReply[fieldOffset : fieldOffset+2])
		case "U32":
			reply[f.name] = binary.BigEndian.Uint32(modbusReply[fieldOffset : fieldOffset+4])
		case "I16":
			reply[f.name] = int16(binary.BigEndian.Uint16(modbusReply[fieldOffset : fieldOffset+2]))
		default:
		}
	}

	return reply, nil

}
