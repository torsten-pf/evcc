package plugin

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/plugin/pipeline"
)

// https://stackoverflow.com/questions/26545883/how-to-do-one-liner-if-else-statement
// IfThenElse evaluates a condition, if true returns the first parameter otherwise the second
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
    if condition {
        return a
    }
    return b
}

// setFormattedValue formats a message template or returns the value formatted as %v if the message template is empty
// if a pipeline is given, this is processed before 
func setFormattedValue(message, param string, v interface{}, pipeline *pipeline.Pipeline) (string, error) {
	payload := fmt.Sprintf("%v", v)
	// check if message contains %s if a pipeline is present
	if (pipeline != nil) && (message != "") && (!strings.Contains(message, "%s")) {
		return "", errors.New("payload must use %s as placeholder if pipeline is used: payload = '" + message + "'")
	}
	if pipeline != nil {
		processed_payload, err := pipeline.Process([]byte(payload))
		if err != nil {
			return "", err
		}
		
		payload = string(processed_payload)		
	}
	
	if message != "" {
		var err error
		payload, err = util.ReplaceFormatted(message, map[string]interface{}{
			param: IfThenElse(pipeline != nil, payload, v),
		})
		if err != nil {
			return "", err
		}
	}

	return payload, nil
}

// knownErrors maps string responses to known error codes
func knownErrors(b []byte) error {
	switch string(b) {
	case "ErrAsleep":
		return api.ErrAsleep
	case "ErrMustRetry":
		return api.ErrMustRetry
	case "ErrNotAvailable":
		return api.ErrNotAvailable
	default:
		return nil
	}
}

func contextLogger(ctx context.Context, log *util.Logger) *util.Logger {
	if ctx != nil {
		if l, ok := ctx.Value(util.CtxLogger).(*util.Logger); ok {
			log = l
		}
	}

	return log
}

// parseFloat rejects NaN and Inf values
func parseFloat(payload string) (float64, error) {
	f, err := strconv.ParseFloat(payload, 64)
	if err == nil && (math.IsNaN(f) || math.IsInf(f, 0)) {
		return 0, fmt.Errorf("invalid float value: %s", payload)
	}
	return f, err
}
