package sensor

import ()

type Payload struct {
	Data string
}

// TODO should validate on init (how to inject into struct init)
// TODO should return bad validation string
func (p *Payload) Validate() error {
	//fields := strings.Split(si.data, ",")
	//if len(s)
	return nil
}

// TODO implement
func (p *Payload) Save() error {
	p.Validate()
	return nil
}
