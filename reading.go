package sensor

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Reading struct {
	DeviceID string
	LogDate  time.Time
	TempF    int
}

// TODO this factory (data) is implementation specific, probably belongs somewhere else
func NewReading(data string) (*Reading, error) {
	// TODO should validation be somewhere else?
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

	return &Reading{
		DeviceID: fields[0],
		LogDate:  t,
		TempF:    temp_f,
	}, nil
}

// TODO also implementation specific, belongs elsewhere
func (p *Reading) Save(*sql.DB) error {
	//age := 21
	//rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	return nil
}
