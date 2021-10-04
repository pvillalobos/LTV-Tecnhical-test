package controller

import (
	"fmt"
	"time"

	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	Utils "github.com/bayronaz/LTV-Tecnhical-test/Helpers"
	service "github.com/bayronaz/LTV-Tecnhical-test/Service"
)

type SongController interface {
	GetReleases() []entity.OutputResponse
	AddSongs(entity.OutputResponse)
	AddNotFoundDates(time.Time)
	getGroupedNotFoundDates() map[string][]string //[]time.Time
	ExistNotFoundDates() bool
	GetDataForNotFoundDates() (bool, error)
}

type controller struct {
	service       service.SongService
	DatesNotFound []time.Time
}

func New(service service.SongService) SongController {
	return &controller{
		service: service,
	}
}

func (c *controller) GetReleases() []entity.OutputResponse {
	return c.service.GetReleases()
}

func (c *controller) AddSongs(entity entity.OutputResponse) {
	c.service.AddSongs(entity)
}

func (c *controller) AddNotFoundDates(date time.Time) {
	c.DatesNotFound = append(c.DatesNotFound, date)
}

func (c *controller) ExistNotFoundDates() bool {
	return len(c.DatesNotFound) > 0
}

//Return grouped in months the not found dates in caché
func (c *controller) getGroupedNotFoundDates() map[string][]string {
	datesGrouped := make(map[string][]string)

	//groups dates in months to evaluate if consume API should be by day or month
	for _, date := range c.DatesNotFound {
		month := date.Format(Utils.Parse_Layout_MM)
		dayInString := string(date.Format(Utils.Parse_Layout))

		datesGrouped[string(month)] = append(datesGrouped[string(month)], dayInString)
	}
	return datesGrouped
}

//This method will make request to songs repository and handle errors
func (c *controller) GetDataForNotFoundDates() (bool, error) {

	for _, data := range c.getGroupedNotFoundDates() {
		if len(data) < 25 {
			for _, date := range data {
				fmt.Println("Fecha simple: ", date)
			}
		} else {
			fmt.Println("Fecha mes: ", data)
		}

	}
	return true, nil
}
