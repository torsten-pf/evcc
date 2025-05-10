package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/evcc-io/evcc/plugin/pipeline"
)

func SetFormattedValueRunner(payload string, jq string, param string, v interface{}) (string, error) {
	if (jq != "") {
		pipeline, err := new(pipeline.Pipeline).WithJq(jq)
		if err != nil {
			return "", err
		}
		return setFormattedValue(payload, param, v, pipeline)
	}
	return setFormattedValue(payload, param, v, nil)
}

func TestSetFormattedValueIntNoPipeline(t *testing.T) {
	{
		payload, err := SetFormattedValueRunner("${var:%d}", "", "var", 1)
		assert.NoError(t, err)
		assert.Equal(t, "1", payload)
	}
}

func TestSetFormattedValueAdd(t *testing.T) {
	{
		payload, err := SetFormattedValueRunner("${var:%s}", ". + 1", "var", 1)
		assert.NoError(t, err)
		assert.Equal(t, "2", payload)
	}
}

func TestSetFormattedValueMul(t *testing.T) {
	{
		payload, err := SetFormattedValueRunner("${var:%s}", ". * 10", "var", 1)
		assert.NoError(t, err)
		assert.Equal(t, "10", payload)
	}
}

func TestSetFormattedValueNoPayloadStr(t *testing.T) {
	{
		payload, err := SetFormattedValueRunner("", ". * 10", "var", 1)
		assert.NoError(t, err)
		assert.Equal(t, "10", payload)
	}
}

func TestSetFormattedValueJson(t *testing.T) {
	{
		payload, err := SetFormattedValueRunner("{\"value\": ${var:%s}}", ". + 1", "var", 1)
		assert.NoError(t, err)
		assert.Equal(t, "{\"value\": 2}", payload)
	}
}

// Error: do not use other types than string in payload formatter
func TestSetFormattedValueErr(t *testing.T) {
	{
		_, err := SetFormattedValueRunner("${var:%d}", ". * 10", "var", 1)
		assert.Error(t, err)				
	}
}