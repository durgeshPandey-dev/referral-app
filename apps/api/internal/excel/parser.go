package excel

import (
	"fmt"
	"strings"

	"github.com/durgeshPandey-dev/referral/apps/api/internal/models"

	"github.com/xuri/excelize/v2"
)

func Parse(filePath string) ([]models.HRContact, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows found")
	}

	var contacts []models.HRContact

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		if len(row) == 0 {
			continue
		}

		var contact models.HRContact

		if len(row) > 0 {
			contact.Email = strings.TrimSpace(row[0])
		}
		if len(row) > 1 {
			contact.Name = strings.TrimSpace(row[1])
		}
		if len(row) > 2 {
			contact.CompanyName = strings.TrimSpace(row[2])
		}

		if contact.Email == "" {
			continue
		}

		contacts = append(contacts, contact)
	}

	if len(contacts) == 0 {
		return nil, fmt.Errorf("no valid contacts found")
	}

	return contacts, nil
}
