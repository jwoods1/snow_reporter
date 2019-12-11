package actions

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"snow_watch/models"
	"strconv"

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

	//dailyReports := &models.DailyReports{}
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
	stations := models.Stations{}
	models.DB.All(&stations)
	fUrl := "https://wcc.sc.egov.usda.gov/reportGenerator/view_csv/customSingleStationReport/daily/"
	eUrl := ":ID:SNTL%7Cid=%22%22%7Cname/-7,0/SNWD::value,SNWD::delta"
	for _, s := range stations {
		url := fUrl + string(s.StationID) + eUrl
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		records := readCsv(resp.Body)
		for _, r := range records[1:] {
			fmt.Println(r)
		}

	}
}

func getAllStations(tx *pop.Connection) {
	url := "https://wcc.sc.egov.usda.gov/reportGenerator/view_csv/customMultipleStationReport/daily/start_of_period/state=%22ID%22%20AND%20network=%22SNTLT%22,%22SNTL%22%20AND%20element=%22SNWD%22%20AND%20outServiceDate=%222100-01-01%22%7Cname/0,0/name,stationId,elevation,latitude,longitude,state.code?fitToScreen=false"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	records := readCsv(resp.Body)
	for _, r := range records[1:] {

		s := &models.Station{}

		s.Name = r[0]
		id, err := strconv.Atoi(r[1])
		if err != nil {
			log.Fatalln(err)
		}
		s.StationID = id
		el, err := strconv.Atoi(r[2])
		if err != nil {
			log.Fatalln(err)
		}
		s.Elevation = el
		la, err := strconv.ParseFloat(r[3], 64)
		s.Lat = la
		lo, err := strconv.ParseFloat(r[4], 64)
		s.Long = lo
		s.State = r[5]

		tx.Create(s)
	}
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
