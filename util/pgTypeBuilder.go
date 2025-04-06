package util

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// ToPgNumeric конвертирует int64 в pgtype.int8
func ToPgInt(value int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64:   value,
		Valid: true,
	}
}

// ToPgBool конвертирует bool в pgtype.Bool
func ToPgBool(value bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  value,
		Valid: true,
	}
}

// ToPgText конвертирует string в pgtype.Text
func ToPgText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}