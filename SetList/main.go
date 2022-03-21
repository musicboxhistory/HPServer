package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

type SetListData struct {
	Num      string `json:"num"`
	Music    string `json:"music"`
	Property string `json:"property"`
	Member   string `json:"member"`
}

var liveList map[string]map[string][]SetListData

func setLiveData() error {
	// Init Data
	group := []string{"angerme", "berryz", "beyonds", "camellia", "country", "cute", "hello", "juice", "magnolia", "morning", "ochanorma", "trainee"}
	liveList = map[string]map[string][]SetListData{}

	// Set Group Live Data
	for _, value := range group {
		fileName := fmt.Sprintf("./xlsx/%s.xlsx", value)
		excel, err := xlsx.OpenFile(fileName)
		if err != nil {
			fmt.Println(err)
			return err
		}

		setLists := map[string][]SetListData{}
		for _, sheet := range excel.Sheets {
			setLists[sheet.Name] = make([]SetListData, 0)
			for idx := 0; idx < sheet.MaxRow; idx++ {
				var setList SetListData
				cell, _ := sheet.Cell(idx, 0)
				setList.Num = cell.Value
				cell, _ = sheet.Cell(idx, 1)
				setList.Music = cell.Value
				cell, _ = sheet.Cell(idx, 2)
				setList.Property = cell.Value
				cell, _ = sheet.Cell(idx, 3)
				setList.Member = cell.Value
				setLists[sheet.Name] = append(setLists[sheet.Name], setList)
			}
		}
		liveList[value] = setLists
	}

	return nil
}

func groupListPrint(c *gin.Context) {
	group := []string{"angerme", "berryz", "beyonds", "camellia", "country", "cute", "hello", "juice", "magnolia", "morning", "ochanorma", "trainee"}
	c.JSON(http.StatusOK, group)
}

func setListPrint(c *gin.Context) {
	// Init Data
	emptyData := make([]SetListData, 0)

	// Get Path Info
	group := c.Param("group")
	live := c.Param("live")

	// Value Check
	if liveList[group][live] == nil {
		fmt.Printf("Not Group or Not Live\n")
		c.JSON(http.StatusNotFound, emptyData)
		return
	}

	c.JSON(http.StatusOK, liveList[group][live])
}

func liveListPrint(c *gin.Context) {
	// Init Data
	liveListData := make([]string, 0)

	// Get Path Info
	group := c.Param("group")

	// Value Check
	if liveList[group] == nil {
		fmt.Printf("Not Group\n")
		c.JSON(http.StatusNotFound, liveListData)
		return
	}

	// Set Live List
	for key, _ := range liveList[group] {
		liveListData = append(liveListData, key)
	}

	c.JSON(http.StatusOK, liveListData)
}

func main() {
	// Set Live Data
	err := setLiveData()
	if err != nil {
		fmt.Println("setLiveData NG")
		return
	}

	// HTTP Server Start
	r := gin.Default()
	r.GET("/HelloProject/SetList", groupListPrint)
	r.GET("/HelloProject/SetList/:group", liveListPrint)
	r.GET("/HelloProject/SetList/:group/:live", setListPrint)
	r.Run(":80")
}
