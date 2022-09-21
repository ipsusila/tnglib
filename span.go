package tnglib

import (
	"encoding/json"
	"time"
)

// Span stores duration for json decode/encode
type Span struct {
	time.Duration
}

// TimeSpan convert duration to Span
func TimeSpan(d time.Duration) Span {
	return Span{Duration: d}
}

// MustTimeSpan convert string to Span
func MustTimeSpan(s string) Span {
	dur, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return Span{Duration: dur}
}

// MarshalJSON converts duration to string
func (d Span) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON convert duration stream to time.Duration
func (d *Span) UnmarshalJSON(data []byte) error {
	// Try unmarshal to number
	var fv float64
	if err := json.Unmarshal(data, &fv); err == nil {
		d.Duration = time.Duration(fv)
		return nil
	}

	// Try unmarshal to string
	var sv string
	if err := json.Unmarshal(data, &sv); err != nil {
		return err
	}
	dur, err := time.ParseDuration(sv)
	if err != nil {
		return err
	}
	d.Duration = dur
	return nil
}
