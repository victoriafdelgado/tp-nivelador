package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type Bet struct {
	AgencyId    int
	Name        string
	Surname     string
	Id          int
	DateOfBirth string
	Number      int
}

func ParseBetFromString(betString string) (Bet, error) {

	fields := strings.Split(betString, ",")

	AgencyId, err := strconv.Atoi(fields[0])
	if err != nil {
		return Bet{}, fmt.Errorf("dni inválido: %w", err)
	}
	Name := fields[1]
	Surname := fields[2]
	Id, err := strconv.Atoi(fields[3])
	if err != nil {
		return Bet{}, fmt.Errorf("dni inválido: %w", err)
	}
	DateOfBirth := fields[4]
	Number, err := strconv.Atoi(fields[5])
	if err != nil {
		return Bet{}, fmt.Errorf("numero inválido: %w", err)
	}

	bet := Bet{
		AgencyId:    AgencyId,
		Name:        Name,
		Surname:     Surname,
		Id:          Id,
		DateOfBirth: DateOfBirth,
		Number:      Number,
	}

	return bet, nil
}
