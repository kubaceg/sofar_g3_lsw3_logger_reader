package ports

import "time"

type CommunicationPort interface {
	Open() error
	Read(buffer []byte) (int, error)
	Write(payload []byte) (int, error)
	Close() error
	SetWriteDeadline(deadline time.Time) error
	SetReadDeadline(deadline time.Time) error
}
