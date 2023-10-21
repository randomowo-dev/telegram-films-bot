package http

import (
	"encoding/json"
	"strings"
)

type FilmProductionStatus int8

const (
	Filming FilmProductionStatus = iota
	PreProduction
	Completed
	Announced
	PostProduction
	UnknownStatus = -1
)

func (s *FilmProductionStatus) String() string {
	switch *s {
	case Filming:
		return "FILMING"
	case PreProduction:
		return "PRE_PRODUCTION"
	case Completed:
		return "COMPLETED"
	case Announced:
		return "ANNOUNCED"
	case PostProduction:
		return "POST_PRODUCTION"
	}

	return "UNKNOWN"
}

func parseFilmProductionStatus(status string) FilmProductionStatus {
	switch strings.ToUpper(status) {
	case "FILMING":
		return Filming
	case "PRE_PRODUCTION":
		return PreProduction
	case "COMPLETED":
		return Completed
	case "ANNOUNCED":
		return Announced
	case "POST_PRODUCTION":
		return PostProduction
	}
	return UnknownStatus
}

func (s *FilmProductionStatus) UnmarshalJSON(data []byte) error {
	var productionStatus string
	if err := json.Unmarshal(data, &productionStatus); err != nil {
		return err
	}

	*s = parseFilmProductionStatus(productionStatus)
	return nil
}

type FilmType int8

const (
	Film FilmType = iota
	Video
	TvSeries
	MiniSeries
	TvShow
	UnknownType = -1
)

func (t *FilmType) String() string {
	switch *t {
	case Film:
		return "FILM"
	case Video:
		return "VIDEO"
	case TvSeries:
		return "TV_SERIES"
	case MiniSeries:
		return "MINI_SERIES"
	case TvShow:
		return "TV_SHOW"
	}

	return "UNKNOWN"
}

func parseFilmType(ft string) FilmType {
	switch strings.ToUpper(ft) {
	case "FILM":
		return Film
	case "VIDEO":
		return Video
	case "TV_SERIES":
		return TvSeries
	case "MINI_SERIES":
		return MiniSeries
	case "TV_SHOW":
		return TvShow
	}
	return UnknownType
}

func (t *FilmType) UnmarshalJSON(data []byte) error {
	var productionStatus string
	if err := json.Unmarshal(data, &productionStatus); err != nil {
		return err
	}

	*t = parseFilmType(productionStatus)
	return nil
}

type Countries []string

func (t *Countries) UnmarshalJSON(data []byte) error {
	var countries []map[string]string
	if err := json.Unmarshal(data, &countries); err != nil {
		return err
	}

	*t = make([]string, 0, len(countries))
	for _, c := range countries {
		if val, ok := c["country"]; ok {
			*t = append(*t, val)
		}
	}

	return nil
}

type Genre int8

const (
	Anime Genre = iota
	Biography
	Action
	Western
	War
	Detective
	ForChild
	Documentary
	Drama
	Historical
	Comedy
	Short
	Criminal
	Melodrama
	Music
	Cartoon
	Musical
	Adventure
	ForFamily
	Sport
	Thriller
	Horror
	Fiction
	Noir
	Fantasy
	UnknownGenre = -1
)

func (g Genre) String() string {
	switch g {
	case Anime:
		return "аниме"
	case Biography:
		return "биография"
	case Action:
		return "боевик"
	case Western:
		return "вестерн"
	case War:
		return "военный"
	case Detective:
		return "детектив"
	case ForChild:
		return "детский"
	case Documentary:
		return "документальный"
	case Drama:
		return "драма"
	case Historical:
		return "исторический"
	case Comedy:
		return "комедия"
	case Short:
		return "короткометражка"
	case Criminal:
		return "криминал"
	case Melodrama:
		return "мелодрама"
	case Music:
		return "музыка"
	case Cartoon:
		return "мультфильм"
	case Musical:
		return "мюзикл"
	case Adventure:
		return "приключения"
	case ForFamily:
		return "cемейный"
	case Sport:
		return "спорт"
	case Thriller:
		return "триллер"
	case Horror:
		return "хоррор"
	case Fiction:
		return "фантастика"
	case Noir:
		return "фэнтези"
	case Fantasy:
		return "фильм-нуар"
	}
	return "неизвестно"
}

func parseGenre(genre string) Genre {
	switch strings.ToLower(genre) {
	case "аниме":
		return Anime
	case "биография":
		return Biography
	case "боевик":
		return Action
	case "вестерн":
		return Western
	case "военный":
		return War
	case "детектив":
		return Detective
	case "детский":
		return ForChild
	case "документальный":
		return Documentary
	case "драма":
		return Drama
	case "исторический":
		return Historical
	case "комедия":
		return Comedy
	case "короткометражка":
		return Short
	case "криминал":
		return Criminal
	case "мелодрама":
		return Melodrama
	case "музыка":
		return Music
	case "мультфильм":
		return Cartoon
	case "мюзикл":
		return Musical
	case "приключения":
		return Adventure
	case "cемейный":
		return ForFamily
	case "спорт":
		return Sport
	case "триллер":
		return Thriller
	case "хоррор":
		return Horror
	case "фантастика":
		return Fiction
	case "фильм-нуар":
		return Noir
	case "фэнтези":
		return Fantasy
	}
	return UnknownGenre
}

type Genres []Genre

func (gs *Genres) UnmarshalJSON(data []byte) error {
	var genres []map[string]string
	if err := json.Unmarshal(data, &genres); err != nil {
		return err
	}

	*gs = make([]Genre, 0, len(genres))
	for _, gObj := range genres {
		if g, ok := gObj["genre"]; ok {
			*gs = append(*gs, parseGenre(g))
		}
	}

	return nil
}

type FilmObj struct {
	KinopoiskID             int64                `json:"kinopoiskId"`
	KinopoiskUID            string               `json:"kinopoiskHDId"`
	ImdbID                  string               `json:"imdbId"`
	NameRu                  string               `json:"nameRu"`
	NameEn                  string               `json:"nameEn"`
	NameOriginal            string               `json:"nameOriginal"`
	PosterUrlPreview        string               `json:"posterUrlPreview"`
	RatingGoodReviewPercent float32              `json:"ratingGoodReview"`
	RatingKinopoisk         float32              `json:"ratingKinopoisk"`
	RatingImdb              float32              `json:"ratingImdb"`
	RatingFilmCritics       float32              `json:"ratingFilmCritics"`
	KinopoiskWebUrl         string               `json:"webUrl"`
	Year                    int                  `json:"year"`
	FilmLengthSeconds       int                  `json:"filmLengthSeconds"`
	ShortDescription        string               `json:"shortDescription"`
	ProductionStatus        FilmProductionStatus `json:"productionStatus"`
	Type                    FilmType             `json:"type"`
	RatingAgeLimit          string               `json:"ratingAgeLimits"`
	Countries               Countries            `json:"countries"`
	Genres                  Genres               `json:"genres"`
	IsSerial                bool                 `json:"serial"`
	IsShortFilm             bool                 `json:"shortFilm"`
	IsCompeted              bool                 `json:"competed"`
}

type Episode struct {
	SeasonNumber  int    `json:"seasonNumber"`
	EpisodeNumber int    `json:"episodeNumber"`
	NameRu        string `json:"nameRu"`
	NameEn        string `json:"nameEn"`
	Synopsis      string `json:"synopsis"`
	ReleaseDate   string `json:"releaseDate"`
}

type Season struct {
	Number   int       `json:"number"`
	Episodes []Episode `json:"episodes"`
}

type Seasons struct {
	Total int      `json:"total"`
	Items []Season `json:"items"`
}

type Similar struct {
	Total int       `json:"total"`
	Items []FilmObj `json:"items"`
}

type ExternalSource struct {
	Url      string `json:"url"`
	Platform string `json:"platform"`
	LogoUrl  string `json:"logoUrl"`
}

type ExternalSources struct {
	Total int              `json:"total"`
	Items []ExternalSource `json:"items"`
}

type SearchResult struct {
}

type FilmFilters struct {
	Genres []struct {
		ID    int   `json:"ID"`
		Genre Genre `json:"genre"`
	} `json:"genres"`
	Countries []struct {
		ID      int    `json:"ID"`
		Country string `json:"country"`
	} `json:"countries"`
}

type SearchOrder string

const (
	ByRating  SearchOrder = "RATING"
	ByNumVote SearchOrder = "NUM_VOTE"
	ByYear    SearchOrder = "YEAR"
)

type SearchParams struct {
	// Countries max length = 1
	Countries []int `json:"countries"`
	// Genres max length = 1
	Genres      []string    `json:"genres"`
	Order       SearchOrder `json:"order"`
	FilmTypeStr string      `json:"filmTypeStr"`
	RatingFrom  int         `json:"ratingFrom"`
	RatingTo    int         `json:"ratingTo"`
	YearFrom    int         `json:"yearFrom"`
	YearTo      int         `json:"yearTo"`
	ImdbID      string      `json:"imdbID"`
	Keyword     string      `json:"keyword"`
}

type StaffProfession int8

const (
	Writer = iota
	Operator
	Editor
	Composer
	Producer
	Translator
	Director
	Design
	Actor
	VoiceDirector
	UnknownProfession = -1
)

func (p *StaffProfession) String() string {
	switch *p {
	case Writer:
		return "WRITER"
	case Operator:
		return "OPERATOR"
	case Editor:
		return "EDITOR"
	case Composer:
		return "COMPOSER"
	case Producer:
		return "PRODUCER"
	case Translator:
		return "TRANSLATOR"
	case Director:
		return "DIRECTOR"
	case Design:
		return "DESIGN"
	case Actor:
		return "ACTOR"
	case VoiceDirector:
		return "VOICE_DIRECTOR"
	}

	return "UNKNOWN"
}

func parseStaffProfession(status string) StaffProfession {
	switch strings.ToUpper(status) {
	case "WRITER":
		return Writer
	case "OPERATOR":
		return Operator
	case "EDITOR":
		return Editor
	case "COMPOSER":
		return Composer
	case "PRODUCER", "PRODUCER_USSR":
		return Producer
	case "TRANSLATOR":
		return Translator
	case "DIRECTOR":
		return Director
	case "DESIGN":
		return Design
	case "ACTOR":
		return Actor
	case "VOICE_DIRECTOR":
		return VoiceDirector
	}
	return UnknownProfession
}

func (p *StaffProfession) UnmarshalJSON(data []byte) error {
	var productionStatus string
	if err := json.Unmarshal(data, &productionStatus); err != nil {
		return err
	}

	*p = parseStaffProfession(productionStatus)
	return nil
}

type FilmStaff []struct {
	ID         int             `json:"staffId"`
	NameRu     string          `json:"nameRu"`
	NameEn     string          `json:"nameEn"`
	Profession StaffProfession `json:"professionKey"`
}

type Staff struct {
	ID              int    `json:"ID"`
	KinopoiskWebUrl string `json:"webUrl"`
	NameRu          string `json:"nameRu"`
	NameEn          string `json:"nameEn"`
	Films           []struct {
		ID         int             `json:"filmId"`
		NameRu     string          `json:"nameRu"`
		NameEn     string          `json:"nameEn"`
		Profession StaffProfession `json:"professionKey"`
	} `json:"films"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}
