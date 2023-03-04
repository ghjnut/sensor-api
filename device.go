package sensor

import "time"

type Device struct {
	// TODO look up functions as fields
	ID   string `json:"deviceId"`
	Logs []Log  `json:"logs"`
	// inferred
	AverageTemperature int       `json:"averageTemperature"`
	MostRecentLogDate  time.Time `json:"mostRecentLogDate"`
	TotalAlerts        int       `json:"totalAlerts"`
}

func (d *Device) SetAverageTemperature() int {
	// could lazy load and save value
	total := 0
	for _, l := range d.Logs {
		total += l.TemperatureF
	}
	d.AverageTemperatue = total / len(d.Logs)
	return d.AverageTemperature
}

// return Time.Zero if there are no logs (zero-value for time.Time)
func (d *Device) SetMostRecentLogDate() (t time.Time) {
	for _, l := range d.Logs {
		// work for first element since t is zeroed
		if l.Date.After(t) {
			t = l.Date
		}
	}
	d.MostRecentLogDate = t
	return d.MostRecentLogDate
}

func (d *Device) SetTotalAlerts() (cnt int) {
	for _, l := range d.Logs {
		if l.Alert {
			cnt++
		}
	}
	d.TotalAlerts = cnt
	return d.TotalAlerts
}
