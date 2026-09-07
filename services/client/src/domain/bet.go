package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type Bet struct {
	Name        string
	Surname     string
	Id          int
	DateOfBirth string
	Number      int
}

const fieldCount = 5

func ParseBetFromString(betString string) (Bet, error) {

	fields := strings.Split(betString, ",")
	if len(fields) < fieldCount {
		return Bet{}, fmt.Errorf("formato inválido, se esperaban %d campos y llegaron %d", fieldCount, len(fields))
	}

	Name := fields[0]
	Surname := fields[1]
	Id, err := strconv.Atoi(fields[2])
	if err != nil {
		return Bet{}, fmt.Errorf("dni inválido: %w", err)
	}
	DateOfBirth := fields[3]
	Number, err := strconv.Atoi(fields[4])
	if err != nil {
		return Bet{}, fmt.Errorf("numero inválido: %w", err)
	}

	bet := Bet{
		Name:        Name,
		Surname:     Surname,
		Id:          Id,
		DateOfBirth: DateOfBirth,
		Number:      Number,
	}

	return bet, nil
}
