package sensor

import "time"

type Device struct {
	// TODO look up functions as fields
	ID         string    `json:"deviceId"`
	Logs       []Log     `json:"logs"`
	AvgTemp    int       `json:"averageTemperature"`
	MostRecent time.Time `json:"mostRecentLogDate"`
}

func (d *Device) AverageTemperature() int {
	// could lazy load and save value
	avg := 0
	for _, l := range d.Logs {
		avg += l.TempF
	}
	return avg
}

// return Time.Zero if there are no logs (zero-value for time.Time)
func (d *Device) MostRecentLogDate() (t time.Time) {
	for _, l := range d.Logs {
		// work for first element since t is zeroed
		if l.Date.After(t) {
			t = l.Date
		}
	}
	return t
}
