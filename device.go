package sensor

import "time"

type Device struct {
	// TODO look up functions as fields
	ID   string `json:"deviceId"`
	Logs []Log  `json:"logs"`
	// inferred
	AvgTemp    int       `json:"averageTemperature"`
	MostRecent time.Time `json:"mostRecentLogDate"`
	TotAlerts  int       `json:"totalAlerts"`
}

func (d *Device) AverageTemperature() int {
	// could lazy load and save value
	total := 0
	for _, l := range d.Logs {
		total += l.TempF
	}
	return total / len(d.Logs)
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

func (d *Device) TotalAlerts() (cnt int) {
	for _, l := range d.Logs {
		if l.Alert {
			cnt++
		}
	}
	return cnt
}
