package json_processor_test

import (
	"encoding/json"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/json_processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStream_Watch(t *testing.T) {
	s := json_processor.ProvideJSONStream()

	actual := make([]interface{}, 0)

	go func() {
		for data := range s.Watch() {
			if data.Error != nil {
				continue
			}
			actual = append(actual, data.Data)
		}
	}()

	s.Start("test/ports.json")

	expected := make([]interface{}, 0)

	data := []string{
		`{"name":"Ajman","city":"Ajman","country":"United Arab Emirates","alias":[],"regions":[],"coordinates":[55.5136433,25.4052165],"province":"Ajman","timezone":"Asia/Dubai","unlocs":["AEAJM"],"code":"52000"}`,
		`{"name":"Abu Dhabi","coordinates":[54.37,24.47],"city":"Abu Dhabi","province":"Abu Z¸aby [Abu Dhabi]","country":"United Arab Emirates","alias":[],"regions":[],"timezone":"Asia/Dubai","unlocs":["AEAUH"],"code":"52001"}`,
	}

	for _, d := range data {
		var current interface{}
		err := json.Unmarshal([]byte(d), &current)
		require.NoError(t, err)
		expected = append(expected, current)
	}

	assert.Equal(t, expected, actual)
}
