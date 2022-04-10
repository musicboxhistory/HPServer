package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

type ScheduleListData struct {
	Date        string `json:"date"`
	Prefectures string `json:"prefectures"`
	Venue       string `json:"venue"`
	Start       string `json:"start"`
	Ticket      string `json:"ticket"`
	Word        string `json:"word"`
}

var scheduleList map[string]map[string][]ScheduleListData

func setScheduleData() error {
	// Init Data
	group := []string{"angerme", "beyonds", "camellia", "hello", "juice", "morning", "ochanorma", "trainee"}
	scheduleList = map[string]map[string][]ScheduleListData{}

	// Set Group Live Data
	for _, value := range group {
		fileName := fmt.Sprintf("./xlsx/%s.xlsx", value)
		excel, err := xlsx.OpenFile(fileName)
		if err != nil {
			fmt.Println(err)
			return err
		}

		scheduleLists := map[string][]ScheduleListData{}
		for _, sheet := range excel.Sheets {
			scheduleLists[sheet.Name] = make([]ScheduleListData, 0)
			for idx := 0; idx < sheet.MaxRow; idx++ {
				var scheduleList ScheduleListData
				cell, _ := sheet.Cell(idx, 0)
				scheduleList.Date = cell.Value
				cell, _ = sheet.Cell(idx, 1)
				scheduleList.Prefectures = cell.Value
				cell, _ = sheet.Cell(idx, 2)
				scheduleList.Venue = cell.Value
				cell, _ = sheet.Cell(idx, 3)
				scheduleList.Start = cell.Value
				cell, _ = sheet.Cell(idx, 4)
				scheduleList.Ticket = cell.Value
				cell, _ = sheet.Cell(idx, 5)
				scheduleList.Word = cell.Value
				scheduleLists[sheet.Name] = append(scheduleLists[sheet.Name], scheduleList)
			}
		}
		scheduleList[value] = scheduleLists
	}

	return nil
}

func scheduleListPrint(c *gin.Context) {
	// Init Data
	emptyData := make([]string, 0)

	// Value Check
	if scheduleList == nil {
		fmt.Printf("Not Data\n")
		c.JSON(http.StatusNotFound, emptyData)
		return
	}

	c.JSON(http.StatusOK, scheduleList)
}

func main() {
	// Set Schedule Data
	err := setScheduleData()
	if err != nil {
		fmt.Println("setScheduleData NG")
		return
	}

	// HTTP Server Start
	r := gin.Default()
	r.GET("/HelloProject/Schedule", scheduleListPrint)
	r.Run(":80")
}
