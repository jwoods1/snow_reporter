package grifts

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"snow_watch/models"
	"strconv"
	"time"

	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("reports", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		stations := models.Stations{}
		models.DB.All(&stations)
		fUrl := "https://wcc.sc.egov.usda.gov/reportGenerator/view_csv/customSingleStationReport/daily/"
		eUrl := ":ID:SNTL%7Cid=%22%22%7Cname/-7,0/SNWD::value,SNWD::delta"
		for _, s := range stations {
			sId := strconv.Itoa(s.StationID)
			url := fUrl + sId + eUrl
			resp, err := http.Get(url)
			if err != nil {
				log.Fatalln(err)
			}
			records := readCsv(resp.Body)
			if len(records) > 0 {
				for _, r := range records[1:] {
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
					models.DB.Create(dr)
				}
			}

		}
		return nil
	})
	grift.Add("stations", func(c *grift.Context) error {
		models.DB.TruncateAll()

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

			models.DB.Save(s)
		}
		return nil
	})

})

func readCsv(csv_file io.Reader) [][]string {
	r := csv.NewReader(csv_file)

	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	return records
}
