package service

import entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"

type SongService interface {
	GetReleases() []entity.OutputResponse
	AddSongs(entity.OutputResponse)
}

type songService struct {
	songs []entity.OutputResponse
}

func New() SongService {
	return &songService{}
}

func (service *songService) GetReleases() []entity.OutputResponse {
	return service.songs
}

func (service *songService) AddSongs(release entity.OutputResponse) {
	service.songs = append(service.songs, release)
}
