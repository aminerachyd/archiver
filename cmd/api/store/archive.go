package store

type archive struct {
	metadata archiveMetadata
	Payload  []byte
}

type archiveMetadata struct {
	SizeInBytes int64
	StoredIn    []string
}

func merge(m1, m2 map[string]archiveMetadata) map[string]archiveMetadata {
	if len(m2) == 0 {
		return m1
	}

	result := map[string]archiveMetadata{}

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		if _, exists := result[k]; exists {
			mergedMetadata := archiveMetadata{
				SizeInBytes: v.SizeInBytes,
				StoredIn:    append(result[k].StoredIn, v.StoredIn...),
			}

			result[k] = mergedMetadata
		}
	}

	return result
}

type storageType int

const (
	Azure storageType = 1 << iota
	FileSystem
	TempFileSystem
)

func (s storageType) toString() string {
	switch s {
	case Azure:
		return "Azure"
	case FileSystem:
		return "FileSystem"
	case TempFileSystem:
		return "TempFileSystem"
	default:
		return "Unknown storage"
	}
}
