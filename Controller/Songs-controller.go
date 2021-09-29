package controller

import (
	"time"

	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	service "github.com/bayronaz/LTV-Tecnhical-test/Service"
)

type SongController interface {
	GetReleases() []entity.OutputResponse
	AddSongs(entity.OutputResponse)
	AddNotFoundDates(time.Time)
	GetNotFoundDates() []time.Time
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

func (c *controller) GetNotFoundDates() []time.Time {
	return c.DatesNotFound
}
