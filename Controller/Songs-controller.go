package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	. "github.com/ahmetb/go-linq/v3"
	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	Utils "github.com/bayronaz/LTV-Tecnhical-test/Helpers"
	service "github.com/bayronaz/LTV-Tecnhical-test/Service"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

type SongController interface {
	getReleases() []entity.OutputResponse
	addSongs(entity.OutputResponse)
	buildResponse([]entity.SongsRepositoryAnswer)
	addNotFoundDates(time.Time)
	getGroupedNotFoundDates() map[string][]string
	existNotFoundDates() bool
	getDataForNotFoundDates(ctx *gin.Context)
	saveDataInCache(string, string, string)
	ProcessReleasesRequest(*gin.Context, time.Time, time.Time, string)
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

// Gets All info matching the request to respont the client
func (c *controller) ProcessReleasesRequest(ctx *gin.Context, from time.Time, until time.Time, artist string) {

	for rd := Utils.RangeDate(from, until); ; {
		date := rd()
		//if func RangeDate return no dates, breaks cycle
		if date.IsZero() {
			break
		}
		//Lets look for what we have already store in caché
		songs, found := Utils.Cache.Get(date.Format(Utils.Parse_Layout))
		if !found {
			c.addNotFoundDates(date)
		} else {
			c.buildResponse(songs.([]entity.SongsRepositoryAnswer))
		}
	}

	//Check if there is missing dates to consume API
	if c.existNotFoundDates() {
		c.getDataForNotFoundDates(ctx)
		c.buildResponse(nil)
	}
	ctx.IndentedJSON(http.StatusOK, c.getReleases())
}

//Fill API_PreResponse with all data that should be answer, but it has to be transform into OutputResponse first
func (c *controller) buildResponse(data []entity.SongsRepositoryAnswer) {
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
			//c.BuildResponse(songs.([]entity.SongsRepositoryAnswer))
			if c.Artist != "" {
				for _, song := range songs.([]entity.SongsRepositoryAnswer) {
					if song.Artist == c.Artist {
						c.API_PreResponse = append(c.API_PreResponse, song)
					}
				}
			} else {
				c.API_PreResponse = append(c.API_PreResponse, songs.([]entity.SongsRepositoryAnswer)...)
			}
		}
	}
}

//Return all data in API_PreResponse as OutputResponse array
func (c *controller) getReleases() []entity.OutputResponse {
	var outputlist = make(map[string]entity.OutputResponse)

	//lets group info in a map with ReleaseAt date and an array of songs
	for _, data := range c.API_PreResponse {
		s := entity.Songs{Artist: data.Artist, Name: data.Name}
		outputlist[data.ReleasedAt] = entity.OutputResponse{ReleasedAt: data.ReleasedAt, Songs: append(outputlist[data.ReleasedAt].Songs, s)}
	}

	//Place all info into service and returned it
	for _, outp := range outputlist {
		c.addSongs(outp)
	}

	return c.service.GetReleases()
}

//Send to the service the songs that should be answered to the request
func (c *controller) addSongs(entity entity.OutputResponse) {
	c.service.AddSongs(entity)
}

//Fill an array with all dates that we do not have in caché
func (c *controller) addNotFoundDates(date time.Time) {
	c.DatesNotFound = append(c.DatesNotFound, date)
}

//Function to validate if there is missing data from dates we do not have in cache
func (c *controller) existNotFoundDates() bool {
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
//It is using - GoRoutines-
func (c *controller) getDataForNotFoundDates(ctx *gin.Context) {
	var wg sync.WaitGroup
	errorChannel := make(chan string, 1)

	for _, data := range c.getGroupedNotFoundDates() {
		if len(data) < 25 {
			for _, date := range data {
				wg.Add(1)
				fmt.Println("Fecha por día: ", date)
				go func(c *controller, ctx *gin.Context, date string, errorChannel chan string) {
					defer wg.Done()
					c.saveDataInCache(Utils.ConsumeSongsRepositoryAPI(date, "daily", ctx, errorChannel), "daily", date)
				}(c, ctx, date, errorChannel)
			}
		} else {
			wg.Add(1)
			date, _ := time.Parse(Utils.Parse_Layout, data[0])
			dateString := string(date.Format(Utils.Parse_Layout_MM))
			fmt.Println("Fecha del mes: ", dateString)
			go func(c *controller, ctx *gin.Context, date string, errorChannel chan string) {
				defer wg.Done()
				c.saveDataInCache(Utils.ConsumeSongsRepositoryAPI(date, "monthly", ctx, errorChannel), "monthly", date)
			}(c, ctx, dateString, errorChannel)
		}
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for errors := range errorChannel {
		fmt.Println(errors)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ErrorResponse{Error: errors})
		return
		//break
	}
	return
}

//Fill caché with info for everyday requested
func (c *controller) saveDataInCache(body string, mode string, date string) {
	if body == "" {
		return
	}

	var jsonInput = []byte(body)
	var dataToStore []entity.SongsRepositoryAnswer

	marshallError := json.Unmarshal(jsonInput, &dataToStore)
	if marshallError != nil {
		log.Fatal("Could not Unmarshall", marshallError)
	}

	if mode == "daily" {
		Utils.Cache.Set(date, dataToStore, cache.DefaultExpiration)
	} else if mode == "monthly" {

		MesActual, _ := time.Parse(Utils.Parse_Layout, date+"-01")
		var sliceTime []time.Time

		//Lets store only dates in the current month, this is why we slice time array because c.DatesNotFound contains all dates not stores in cache
		//it could be from several months
		From(c.DatesNotFound).Where(func(q interface{}) bool {
			return q.(time.Time).Month() == MesActual.Month()
		}).Select(func(q interface{}) interface{} {
			return q.(time.Time)
		}).ToSlice(&sliceTime)

		for _, daysMissing := range sliceTime {
			var songsToSave []entity.SongsRepositoryAnswer
			dateToStore := string(daysMissing.Format(Utils.Parse_Layout))

			//almacenamos solo del mismo mes
			From(dataToStore).Where(func(q interface{}) bool {
				return q.(entity.SongsRepositoryAnswer).ReleasedAt == dateToStore
			}).Select(func(q interface{}) interface{} {
				return q.(entity.SongsRepositoryAnswer)
			}).ToSlice(&songsToSave)

			Utils.Cache.Set(dateToStore, songsToSave, cache.DefaultExpiration)
		}
	}

}
