package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Data is taken from https://www.kaggle.com/datasets/zhongtr0n/country-flag-urls
	file, err := os.Open("flags_iso.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	var (
		baseInsertCommand = "INSERT INTO countries (id, name, flag_url, enable_registration)\nVALUES\n"
		insertValues      []string

		baseDeleteCommand = "DELETE FROM countries WHERE id IN"
		idsToDelete       []string
	)

	// Skipping the header
	for _, row := range records[1:] {
		var (
			countryAlpha2Code  = row[1]
			countryName        = row[0]
			countryFlagUrl     = row[3]
			enableRegistration = false

			safeCountryName = strings.ReplaceAll(countryName, "'", "''") // Some country name contains "'", it needs to be escaped.
		)

		if countryAlpha2Code == "ID" || countryAlpha2Code == "SG" {
			enableRegistration = true
		}

		insertValues = append(insertValues, fmt.Sprintf("('%s', '%s', '%s', %t)", countryAlpha2Code, safeCountryName, countryFlagUrl, enableRegistration))
		idsToDelete = append(idsToDelete, fmt.Sprintf("'%s'", countryAlpha2Code))
	}

	var (
		insertCommand = fmt.Sprintf("%s%s;", baseInsertCommand, strings.Join(insertValues, ",\n"))
		deleteCommand = fmt.Sprintf("%s(%s);", baseDeleteCommand, strings.Join(idsToDelete, ","))
	)

	fmt.Println(insertCommand)
	fmt.Print("\n\n==============================================================\n\n")
	fmt.Println(deleteCommand)
}
