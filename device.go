package sensor

import "time"

type Device struct {
	DeviceID string
	Logs     []Reading
	//averageTemperature int
	//mostRecentLogDate time.Time
}

func (d *Device) AverageTemp() int {
	// could lazy load and save value
	avg := 0
	for _, log := range d.Logs {
		avg += log.TempF
	}
	return avg
}

// return nil if there are no logs
func (d *Device) MostRecentLogDate() (t time.Time) {
	for _, log := range d.Logs {
		if t == nil || log.LogDate.After(t) {
			t = log.LogDate
		}
	}
	return t
}
