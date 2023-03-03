package sensor

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Log struct {
	DeviceID string
	Date     time.Time
	TempF    int
}

// TODO this factory (data) is implementation specific, probably belongs somewhere else
func NewLog(data string) (*Log, error) {
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

	return &Log{
		DeviceID: fields[0],
		Date:     t,
		TempF:    temp_f,
	}, nil
}

// TODO also implementation specific, belongs elsewhere
func (p *Log) Save(db *sql.DB) error {
	sqlStatement := `
INSERT INTO logs (event_date, device_id, temp_farenheit)
VALUES ($1, $2, $3)
RETURNING event_id`
	event_id := 0
	err := db.QueryRow(sqlStatement, p.Date, p.DeviceID, p.TempF).Scan(&event_id)
	if err != nil {
		return err
	}
	fmt.Println("New record ID is:", event_id)
	return nil
}
