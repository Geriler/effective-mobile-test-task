package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func GoogleUUIDToPgxUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func PgxUUIDToGoogleUUID(id pgtype.UUID) (uuid.UUID, error) {
	if !id.Valid {
		return uuid.Nil, fmt.Errorf("invalid uuid")
	}
	return uuid.FromBytes(id.Bytes[:])
}
