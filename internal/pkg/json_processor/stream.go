package json_processor

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func (s Stream) Watch() <-chan Entry {
	return s.stream
}

func (s Stream) Start(path string) {
	defer close(s.stream)

	file, err := os.Open(path)
	if err != nil {
		s.stream <- Entry{Error: fmt.Errorf("open file: %w", err)}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode opening delimiter: %w", err)}
		return
	}

	i := 1
	for decoder.More() {
		_, err := decoder.Token()
		if err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode key %d: %w", i, err)}
			return
		}

		var data interface{}
		if err := decoder.Decode(&data); err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode value %d: %w", i, err)}
			return
		}
		s.stream <- Entry{Data: data}

		i++
	}

	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode closing delimiter: %w", err)}
		return
	}
}