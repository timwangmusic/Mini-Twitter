package util

import "github.com/google/uuid"

func genUUID() uuid.UUID{
	return uuid.New()
}
