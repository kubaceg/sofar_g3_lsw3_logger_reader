package ports

type CommunicationPort interface {
	Open() error
	Read(buffer []byte) (int, error)
	Write(payload []byte) (int, error)
	Close() error
}
