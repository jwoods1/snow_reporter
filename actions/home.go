package actions

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"snow_watch/models"
	"strconv"
	"sync"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	report := models.DailyReport{}
	models.DB.Order("Date Desc").First(&report)
	if timeFix(report.Date) != timeFix(time.Now()) {
		getDailyReport(tx)
	}
	stations := &models.Stations{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	//q := tx.PaginateFromParams(c.Params())
	// Retrieve all DailyReports from the DB
	if err := tx.Eager().All(stations); err != nil {
		return err
	}
	c.Set("stations", stations)
	return c.Render(200, r.HTML("index.html"))
}

func getDailyReport(tx *pop.Connection) {
	models.DB.RawQuery("TRUNCATE daily_reports").Exec()
	stations := models.Stations{}
	tx.All(&stations)
	var wg sync.WaitGroup
	for _, s := range stations {
		wg.Add(1)
		go stationReports(s, tx, &wg)
	}
	wg.Wait()
}
func stationReports(s models.Station, tx *pop.Connection, wg *sync.WaitGroup) {
	fUrl := "https://wcc.sc.egov.usda.gov/reportGenerator/view_csv/customSingleStationReport/daily/"
	eUrl := ":ID:SNTL%7Cid=%22%22%7Cname/-7,0/SNWD::value,SNWD::delta"
	sId := strconv.Itoa(s.StationID)
	url := fUrl + sId + eUrl
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	records := readCsv(resp.Body)

	if len(records) > 0 {
		for _, r := range records[1:] {
			saveRecords(s, r, tx)
		}
	}
	wg.Done()
}
func saveRecords(s models.Station, r []string, tx *pop.Connection) {
	dr := &models.DailyReport{}
	date, err := time.Parse("2006-01-02", r[0])
	if err != nil {
		log.Fatalln(err)
	}
	dr.Date = date
	dr.StationID = s.ID
	depth, _ := strconv.ParseFloat(r[1], 64)
	dr.SnowDepth = depth
	change, _ := strconv.ParseFloat(r[2], 64)
	dr.SnowChange = change
	tx.Create(dr)
}

func readCsv(csv_file io.Reader) [][]string {
	r := csv.NewReader(csv_file)

	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	return records
}
