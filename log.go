package sensor

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Log struct {
	DeviceID     string    `json:"deviceId,omitempty"`
	Date         time.Time `json:"logDate"`
	TemperatureF int       `json:"temperature"`
	// inferred
	Alert bool `json:"alert"`
}

// TODO this factory (data string) is implementation specific, probably belongs somewhere else
func NewLog(data string) (*Log, error) {
	fields := strings.Split(data, "|")
	if len(fields) != 3 {
		return nil, errors.New("bad payload")
	}

	//2022-03-28T18:55:35.433Z
	t, err := time.Parse(time.RFC3339, fields[1])
	if err != nil {
		return nil, err
	}

	//`TempF = (2*TempC) + 30`
	temp_c, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, err
	}
	temp_f := 2*temp_c + 30

	l := &Log{
		DeviceID:     fields[0],
		Date:         t,
		TemperatureF: temp_f,
	}
	_ = l.SetAlert()

	return l, nil
}

func (d *Log) SetAlert() bool {
	if d.TemperatureF > 32 {
		d.Alert = true
	}
	return d.Alert
}
