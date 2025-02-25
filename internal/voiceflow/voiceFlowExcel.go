package voiceflow

import (
	"basicScraper/internal/schemas"
	"fmt"
	"log"
	"time"

	"github.com/tealeg/xlsx"
)

func ExtractQuestions(filename string) ([]string, error) {
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if len(file.Sheets) != 1 {
		log.Fatal("Invalid Excel spreadsheet")
	}

	questions := make([]string, 0)
	for _, row := range file.Sheets[0].Rows {
		if len(row.Cells) > 0 {
			questions = append(questions, row.Cells[0].String())
		}
	}
	return questions, nil
}

func MakeResultExcelDoc(rowMap []schemas.VoiceFlowExcelRow) error {

	excelFile := xlsx.NewFile()
	currentTime := time.Now()
	currentDateString := currentTime.Format("2006-01-02") // YYYY-MM-DD
	sheet, err := excelFile.AddSheet("VoiceFlowChat-" + currentDateString)
	if err != nil {
		fmt.Printf("Error creating sheet: %s\n", err)
		return err
	}

	headers := []string{"Question", "Answer"}

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

	for _, val := range rowMap {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.SetValue(val.Question)
		style := styleCell()
		cell.SetStyle(style)

		cell = row.AddCell()
		cell.SetValue(val.Response)
		style = styleCell()
		cell.SetStyle(style)
	}

	// Save the XLSX file
	err = excelFile.Save("VoiceFlowChat" + currentDateString + ".xlsx")
	if err != nil {
		fmt.Println("Error saving file:", err)
		return err
	}

	log.Println("Successfully made: ", "VoiceFlowChat"+currentDateString+".xlsx")

	return nil
}

func styleCell() *xlsx.Style {
	style := xlsx.NewStyle()
	style.Font = *xlsx.NewFont(11, "Times New Roman")
	return style
}
