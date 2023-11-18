package store

type archive struct {
	metadata archiveMetadata
	Payload  []byte
}

type archiveMetadata struct {
	Name        string
	SizeInBytes int64
}
