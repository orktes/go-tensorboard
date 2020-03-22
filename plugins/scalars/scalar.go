package scalars

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// ScalarValue struct contains a single scalar value
type ScalarValue struct {
	//wall_time, step, value
	WallTime time.Time
	Step     int
	Value    float64
}

// MarshalJSON marshal to Tensorboard scalar format
func (sv *ScalarValue) MarshalJSON() ([]byte, error) {
	b := bytes.NewBuffer(nil)

	b.WriteRune('[')

	b.WriteString(fmt.Sprintf("%f", float64(sv.WallTime.UnixNano())/float64(1000000000)))

	b.WriteRune(',')

	b.WriteString(strconv.Itoa(sv.Step))

	b.WriteRune(',')

	b.WriteString(fmt.Sprintf("%f", sv.Value))

	b.WriteRune(']')

	return b.Bytes(), nil
}
