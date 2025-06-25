package gogacon

// Serializer defines interface for configuration serialization
// Implementations should handle marshaling/unmarshaling of config formats
type Serializer interface {
	//Marshal converts configuration to byte representation
	Marshal() ([]byte, error)
	//Unmarshal parses configuration from byte data
	Unmarshal(data []byte) error
}
