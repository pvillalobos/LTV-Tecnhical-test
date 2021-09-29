package entity

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

//Strut from SongsRepository API anwser
type SongsRepositoryAnswer struct {
	SongID     string `json:"song_id"`
	ReleasedAt string `json:"released_at"`
	Duration   string `json:"duration"`
	Artist     string `json:"artist"`
	Name       string `json:"name"`
	Stats      struct {
		LastPlayedAt int64 `json:"last_played_at"`
		TimesPlayed  int   `json:"times_played"`
		GlobalRank   int   `json:"global_rank"`
	} `json:"stats,omitempty"`
}

//Final answer struts
type OutputResponse struct {
	ReleasedAt string  `json:"released_at"`
	Songs      []Songs `json:"songs"`
}
type Songs struct {
	Artist string `json:"artist"`
	Name   string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
