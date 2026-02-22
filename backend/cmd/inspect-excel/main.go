package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("/Users/josebeltran/Documents/GitHub/customermx/eventos.xlsx")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	sheet := "General"
	// Print first 3 rows, columns A through BJ
	for row := 1; row <= 3; row++ {
		fmt.Printf("=== Row %d ===\n", row)
		for col := 1; col <= 62; col++ {
			cell := columnName(col) + fmt.Sprintf("%d", row)
			val, _ := f.GetCellValue(sheet, cell)
			if val != "" {
				fmt.Printf("  %s (col %2d): [%s]\n", cell, col, val)
			}
		}
	}
}

func columnName(index int) string {
	name := ""
	for index > 0 {
		index--
		name = string(rune('A'+(index%26))) + name
		index /= 26
	}
	return name
}
