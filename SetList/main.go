package main

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

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
var mutex sync.RWMutex

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
				cell := sheet.Cell(idx, 0)
				setList.Num = cell.Value
				cell = sheet.Cell(idx, 1)
				setList.Music = cell.Value
				cell = sheet.Cell(idx, 2)
				setList.Property = cell.Value
				cell = sheet.Cell(idx, 3)
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

	// Mutex Lock
	mutex.RLock()
	defer mutex.RUnlock()

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

	// sort
	sort.Slice(liveListData, func(i, j int) bool { return liveListData[i] < liveListData[j] })

	c.JSON(http.StatusOK, liveListData)
}

func artistPrint(c *gin.Context) {
	c.HTML(http.StatusOK, "artist.html", gin.H{})
}

func livePrint(c *gin.Context) {

	// Mutex Lock
	mutex.RLock()
	defer mutex.RUnlock()

	// Init Data
	liveListData := make([]string, 0)
	artist := c.PostForm("artist")

	// Value Check
	if liveList[artist] == nil {
		fmt.Printf("Not Group\n")
		c.HTML(http.StatusNotFound, "notfound.html", gin.H{})
		return
	}

	// Set Live List
	for key, _ := range liveList[artist] {
		liveListData = append(liveListData, key)
	}

	// sort
	sort.Slice(liveListData, func(i, j int) bool { return liveListData[i] < liveListData[j] })

	c.HTML(http.StatusOK, "live.html", gin.H{"liveList": liveListData, "artist": artist})
}

func setPrint(c *gin.Context) {

	// Mutex Lock
	mutex.RLock()
	defer mutex.RUnlock()

	// Get Post Info
	artist := c.PostForm("artist")
	live := c.PostForm("live")

	fmt.Printf("artist = %s¥n", artist)
	fmt.Printf("live = %s¥n", live)

	// Value Check
	if liveList[artist][live] == nil {
		fmt.Printf("Not Group or Not Live\n")
		c.HTML(http.StatusNotFound, "notfound.html", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "set.html", gin.H{"liveList": liveList[artist][live], "artist": artist, "live": live})
}

func main() {
	// Set Live Data
	go func() {
		for {
			mutex.Lock()
			err := setLiveData()
			if err != nil {
				fmt.Println("setLiveData NG")
				mutex.Unlock()
				return
			}
			mutex.Unlock()
			time.Sleep(3600 * time.Second)
		}
	}()

	// HTTP Server Start
	r := gin.Default()
	r.LoadHTMLGlob("./html/*.html")
	r.Static("/assets", "./assets")
	r.GET("/HelloProject/Concert", artistPrint)
	r.POST("/HelloProject/Live", livePrint)
	r.POST("/HelloProject/Set", setPrint)
	r.GET("/HelloProject/SetList", groupListPrint)
	r.GET("/HelloProject/SetList/:group", liveListPrint)
	r.GET("/HelloProject/SetList/:group/:live", setListPrint)
	r.Run(":80")
}
