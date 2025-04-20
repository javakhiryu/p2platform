package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func DerefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func DerefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func DerefBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func GetValidTime(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
