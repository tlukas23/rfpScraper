package rfp

import (
	"basicScraper/internal/httpReq"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type RFP struct {
	Hyperlink  string
	Name       string
	Agency     string
	Location   string
	DateIssued string
}

func MakeRfpSpreadsheet() (string, error) {
	rsp, _ := httpReq.MakeHttpRequest(http.MethodGet, "https://www.findrfp.com/service/search.aspx?t=FE&s=background+check&x=0&y=0", nil, nil)

	re := regexp.MustCompile(`<a href='([^']+)'>.*?<font color=".*?">(.*?)</font></a></td><td>(.*?)</td><td>\s*(.*?)\s*</td><td>(\d{2}/\d{2}/\d{4})</td>`)
	matches := re.FindAllStringSubmatch(string(rsp), -1)

	rfpArray := make([]RFP, 0)
	for _, match := range matches {
		if len(match) > 4 {
			rfp := RFP{
				Hyperlink:  "https://www.findrfp.com/service/" + match[1],
				Name:       match[2],
				Agency:     match[3],
				Location:   strings.TrimSpace(match[4]),
				DateIssued: match[5],
			}
			rfpArray = append(rfpArray, rfp)
		}
	}

	excelFile := xlsx.NewFile()
	currentTime := time.Now()
	currentDateString := currentTime.Format("2006-01-02") // YYYY-MM-DD
	sheet, err := excelFile.AddSheet("RFPLeads-" + currentDateString)
	if err != nil {
		fmt.Printf("Error creating sheet: %s\n", err)
		return "", err
	}

	headers := []string{"Name", "Agency", "Location", "Issued", "Hyperlink"}

	headerRow := sheet.AddRow()
	for i, header := range headers {
		style := xlsx.NewStyle()
		style.Border.Bottom = "thin"
		style.Font = *xlsx.NewFont(12, "Times New Roman")
		cell := headerRow.AddCell()
		cell.Value = header
		cell.SetStyle(style)
		sheet.SetColWidth(i, i, float64(30)) // Set the width for the first column (0-based index)
	}

	for _, rfp := range rfpArray {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.SetValue(rfp.Name)
		style := styleCell()
		cell.SetStyle(style)

		cell = row.AddCell()
		cell.SetValue(rfp.Agency)
		style = styleCell()
		cell.SetStyle(style)

		cell = row.AddCell()
		cell.SetValue(rfp.Location)
		style = styleCell()
		cell.SetStyle(style)

		cell = row.AddCell()
		cell.SetValue(rfp.DateIssued)
		style = styleCell()
		cell.SetStyle(style)

		cell = row.AddCell()
		cell.SetValue(rfp.Hyperlink)
		style = styleCell()
		cell.SetStyle(style)
	}

	// Save the XLSX file
	err = excelFile.Save("RFPLeads.xlsx")
	if err != nil {
		fmt.Println("Error saving file:", err)
		return "", err
	}
	return "RFPLeads.xlsx", nil
}

func styleCell() *xlsx.Style {
	style := xlsx.NewStyle()
	style.Font = *xlsx.NewFont(11, "Times New Roman")
	return style
}
