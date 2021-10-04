package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	. "github.com/ahmetb/go-linq/v3"
	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	Utils "github.com/bayronaz/LTV-Tecnhical-test/Helpers"
	service "github.com/bayronaz/LTV-Tecnhical-test/Service"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

type SongController interface {
	GetReleases() []entity.OutputResponse
	addSongs(entity.OutputResponse)
	BuildResponse([]entity.SongsRepositoryAnswer)
	AddNotFoundDates(time.Time)
	getGroupedNotFoundDates() map[string][]string //[]time.Time
	ExistNotFoundDates() bool
	GetDataForNotFoundDates(ctx *gin.Context)
	saveDataInCache(string, string, string)
}

type controller struct {
	service         service.SongService
	DatesNotFound   []time.Time
	API_PreResponse []entity.SongsRepositoryAnswer
	Artist          string
}

func New(service service.SongService, artist string) SongController {
	return &controller{
		service: service,
		Artist:  artist,
	}
}

func (c *controller) BuildResponse(data []entity.SongsRepositoryAnswer) {
	if data != nil {
		if c.Artist != "" {
			for _, songs := range data {
				if songs.Artist == c.Artist {
					c.API_PreResponse = append(c.API_PreResponse, songs)
				}
			}
		} else {
			c.API_PreResponse = append(c.API_PreResponse, data...)
		}
	} else {
		//We used DatesNotFound and caché to fill again the API_PreResponse
		for _, date := range c.DatesNotFound {
			songs, _ := Utils.Cache.Get(date.Format(Utils.Parse_Layout))
			c.BuildResponse(songs.([]entity.SongsRepositoryAnswer))
		}
	}
}

func (c *controller) GetReleases() []entity.OutputResponse {
	var outputlist = make(map[string]entity.OutputResponse)

	for _, data := range c.API_PreResponse {
		s := entity.Songs{Artist: data.Artist, Name: data.Name}
		outputlist[data.ReleasedAt] = entity.OutputResponse{ReleasedAt: data.ReleasedAt, Songs: append(outputlist[data.ReleasedAt].Songs, s)}
	}

	//var groupedOutput []entity.OutputResponse
	for _, outp := range outputlist {
		//groupedOutput = append(groupedOutput, outp)
		c.addSongs(outp)
	}

	return c.service.GetReleases()
}

func (c *controller) addSongs(entity entity.OutputResponse) {
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
func (c *controller) GetDataForNotFoundDates(ctx *gin.Context) {

	for _, data := range c.getGroupedNotFoundDates() {
		if len(data) < 25 {
			for _, date := range data {
				fmt.Println("Fecha por día: ", date)
				c.saveDataInCache(Utils.ConsumeSongsRepositoryAPI(date, "daily", ctx), "daily", date)
			}
		} else {
			date, _ := time.Parse(Utils.Parse_Layout, data[0])
			fmt.Println("Fecha del mes: ", string(date.Format(Utils.Parse_Layout_MM)))
			c.saveDataInCache(Utils.ConsumeSongsRepositoryAPI(string(date.Format(Utils.Parse_Layout_MM)), "monthly", ctx), "monthly", "")
		}
	}
}

func (c *controller) saveDataInCache(body string, mode string, date string) {
	var jsonInput = []byte(body)
	var dataToStore []entity.SongsRepositoryAnswer

	marshallError := json.Unmarshal(jsonInput, &dataToStore)
	if marshallError != nil {
		log.Fatal("Could not Unmarshall", marshallError)
	}

	if mode == "daily" {
		Utils.Cache.Set(date, dataToStore, cache.DefaultExpiration)
	} else if mode == "monthly" {
		for _, daysMissing := range c.DatesNotFound {
			var songsToSave []entity.SongsRepositoryAnswer
			date = string(daysMissing.Format(Utils.Parse_Layout))

			From(dataToStore).Where(func(q interface{}) bool {
				return q.(entity.SongsRepositoryAnswer).ReleasedAt == date
			}).Select(func(q interface{}) interface{} {
				return q.(entity.SongsRepositoryAnswer)
			}).ToSlice(&songsToSave)

			Utils.Cache.Set(date, songsToSave, cache.DefaultExpiration)
		}
	}

}
