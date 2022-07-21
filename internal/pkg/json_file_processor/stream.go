package json_file_processor

type Entry struct {
	Error error
	Data  interface{}
}

type Stream struct {
	stream chan Entry
}

func ProvideJSONStream() Stream {
	return Stream{
		stream: make(chan Entry),
	}
}


