package scalars

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Value struct contains a single scalar value
type Value struct {
	//wall_time, step, value
	WallTime time.Time
	Step     int
	Value    float64
}

// MarshalJSON marshal to Tensorboard scalar format
func (sv *Value) MarshalJSON() ([]byte, error) {
	b := bytes.NewBuffer(nil)

	b.WriteRune('[')

	b.WriteString(fmt.Sprintf("%f", float64(sv.WallTime.UnixNano())/float64(1e9)))

	b.WriteRune(',')

	b.WriteString(strconv.Itoa(sv.Step))

	b.WriteRune(',')

	b.WriteString(fmt.Sprintf("%f", sv.Value))

	b.WriteRune(']')

	return b.Bytes(), nil
}
