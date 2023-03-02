package sensor

import ()

type SensorPayload struct {
	Data string
}

// TODO should validate on init (how to inject into struct init)
// TODO should return bad validation string
func (sp *SensorPayload) Validate() error {
	//fields := strings.Split(si.data, ",")
	//if len(s)
	return nil
}

// TODO implement
func (sp *SensorPayload) Save() error {
	sp.Validate()
	return nil
}
