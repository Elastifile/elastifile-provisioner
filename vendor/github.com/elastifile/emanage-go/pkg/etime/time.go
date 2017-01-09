package etime

import "time"

type NilableTime *time.Time

func NewNilableTime(t time.Time) NilableTime {
	return &t
}
