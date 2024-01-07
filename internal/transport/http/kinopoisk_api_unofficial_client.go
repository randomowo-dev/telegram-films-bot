package http

import (
	"encoding/json"
	"fmt"
	"io"
	netHttp "net/http"
	netUrl "net/url"
	"reflect"
	"strconv"

	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	models "github.com/randomowo-dev/telegram-films-bot/internal/models/http"
	"github.com/randomowo-dev/telegram-films-bot/pkg/transport/http/client"
	httpUtils "github.com/randomowo-dev/telegram-films-bot/pkg/utils/http"

	"go.uber.org/zap"
)

type KinopoiskApiUnofficialClient struct {
	client *client.Client
}

type kinopoiskApiUnofficialPath string

var (
	filmsPath kinopoiskApiUnofficialPath = "films"
	staffPath kinopoiskApiUnofficialPath = "staff"
)

func (c *KinopoiskApiUnofficialClient) GetFilmInfoById(filmID int) (*models.FilmObj, error) {
	var (
		err  error
		req  *netHttp.Request
		body []byte
	)
	req, err = c.getBaseRequest(filmsPath, config.KinopoiskApiUnofficialFilmsApiVersion)
	if err != nil {
		return nil, err
	}

	req.URL = req.URL.JoinPath(strconv.Itoa(filmID))

	body, err = c.do(req)
	if err != nil {
		return nil, err
	}

	film := new(models.FilmObj)
	if err = json.Unmarshal(body, film); err != nil {
		c.client.Logger.Error("error on unpack FilmObj", zap.Error(err))
		return nil, err
	}

	return film, nil
}

func (c *KinopoiskApiUnofficialClient) GetFilmSeasons(filmID int) (*models.Seasons, error) {
	var (
		err  error
		req  *netHttp.Request
		body []byte
	)

	req, err = c.getBaseRequest(filmsPath, config.KinopoiskApiUnofficialFilmsApiVersion)
	if err != nil {
		return nil, err
	}

	req.URL = req.URL.JoinPath(strconv.Itoa(filmID), "seasons")

	body, err = c.do(req)
	if err != nil {
		return nil, err
	}

	seasons := new(models.Seasons)
	if err = json.Unmarshal(body, seasons); err != nil {
		c.client.Logger.Error("error on unpack Seasons", zap.Error(err))
		return nil, err
	}

	return seasons, nil
}

func (c *KinopoiskApiUnofficialClient) GetSimilar(filmID int) (*models.Similar, error) {
	var (
		err  error
		req  *netHttp.Request
		body []byte
	)

	req, err = c.getBaseRequest(filmsPath, config.KinopoiskApiUnofficialFilmsApiVersion)
	if err != nil {
		return nil, err
	}

	req.URL = req.URL.JoinPath(strconv.Itoa(filmID), "similars")

	body, err = c.do(req)
	if err != nil {
		return nil, err
	}

	similar := new(models.Similar)
	if err = json.Unmarshal(body, similar); err != nil {
		c.client.Logger.Error("error on unpack Similar", zap.Error(err))
		return nil, err
	}

	return similar, nil
}

func (c *KinopoiskApiUnofficialClient) GetFilmExternalSources(filmID int, page int) (*models.ExternalSources, error) {
	var (
		err  error
		req  *netHttp.Request
		body []byte
	)

	req, err = c.getBaseRequest(filmsPath, config.KinopoiskApiUnofficialFilmsApiVersion)
	if err != nil {
		return nil, err
	}

	req.URL = req.URL.JoinPath(strconv.Itoa(filmID), "external_sources")
	query := req.URL.Query()
	query.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()

	body, err = c.do(req)
	if err != nil {
		return nil, err
	}

	sources := new(models.ExternalSources)
	if err = json.Unmarshal(body, sources); err != nil {
		c.client.Logger.Error("error on unpack ExternalSources", zap.Error(err))
		return nil, err
	}

	return sources, nil
}

func (c *KinopoiskApiUnofficialClient) GetFilmFilters() (*models.FilmFilters, error) {
	var (
		err  error
		req  *netHttp.Request
		body []byte
	)

	req, err = c.getBaseRequest(filmsPath, config.KinopoiskApiUnofficialFilmsApiVersion)
	if err != nil {
		return nil, err
	}

	req.URL = req.URL.JoinPath("filters")

	body, err = c.do(req)
	if err != nil {
		return nil, err
	}

	sources := new(models.FilmFilters)
	if err = json.Unmarshal(body, sources); err != nil {
		c.client.Logger.Error("error on unpack FilmFilters", zap.Error(err))
		return nil, err
	}

	return sources, nil
}

func (c *KinopoiskApiUnofficialClient) Search(searchParams models.SearchParams, page int) (
	*models.SearchResult, error,
) {
	var (
		err  error
		req  *netHttp.Request
		body []byte
	)

	req, err = c.getBaseRequest(filmsPath, config.KinopoiskApiUnofficialFilmsApiVersion)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if err = httpUtils.ParseToQuery(reflect.ValueOf(searchParams), "", &query); err != nil {
		c.client.Logger.Warn("error in parsing search params to query", zap.Error(err))
	}

	query.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()

	body, err = c.do(req)
	if err != nil {
		return nil, err
	}

	result := new(models.SearchResult)
	if err = json.Unmarshal(body, result); err != nil {
		c.client.Logger.Error("error on unpack SearchResult", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (c *KinopoiskApiUnofficialClient) GetStaff(filmID int) (*models.FilmStaff, error) {
	var (
		err  error
		req  *netHttp.Request
		body []byte
	)

	req, err = c.getBaseRequest(staffPath, config.KinopoiskApiUnofficialStaffApiVersion)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("filmId", strconv.Itoa(filmID))
	req.URL.RawQuery = query.Encode()

	body, err = c.do(req)
	if err != nil {
		return nil, err
	}

	filmStaff := new(models.FilmStaff)
	if err = json.Unmarshal(body, &filmStaff); err != nil {
		c.client.Logger.Error("error on unpack FilmStaff", zap.Error(err))
		return nil, err
	}

	return filmStaff, nil
}

func (c *KinopoiskApiUnofficialClient) GetStaffById(staffId int) (*models.Staff, error) {
	var (
		err  error
		req  *netHttp.Request
		body []byte
	)

	req, err = c.getBaseRequest(staffPath, config.KinopoiskApiUnofficialStaffApiVersion)
	if err != nil {
		return nil, err
	}

	req.URL = req.URL.JoinPath(strconv.Itoa(staffId))

	body, err = c.do(req)
	if err != nil {
		return nil, err
	}

	staff := new(models.Staff)
	if err = json.Unmarshal(body, &staff); err != nil {
		c.client.Logger.Error("error on unpack Staff", zap.Error(err))
		return nil, err
	}

	return staff, nil
}

func (c *KinopoiskApiUnofficialClient) getBaseRequest(
	path kinopoiskApiUnofficialPath,
	apiVersion string,
) (*netHttp.Request, error) {
	var (
		url *netUrl.URL
		req *netHttp.Request
		err error
	)

	url, err = netUrl.Parse(config.KinopoiskApiUnofficialUrl)
	if err != nil {
		c.client.Logger.Error(
			"failed to parse kinopoisk api unofficial url",
			zap.String("url", config.KinopoiskApiUnofficialUrl),
			zap.Error(err),
		)
		return nil, err
	}

	url = url.JoinPath("api", apiVersion, string(path))
	req, err = netHttp.NewRequest(netHttp.MethodGet, url.String(), nil)
	if err != nil {
		c.client.Logger.Error(
			"failed to create new request to kinopoisk api unofficial",
			zap.String("method", netHttp.MethodGet),
			zap.String("url", url.String()),
			zap.ByteString("body", nil),
			zap.Error(err),
		)
		return nil, err
	}

	req.Header.Set("X-API-KEY", config.KinopoiskApiUnofficialToken)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *KinopoiskApiUnofficialClient) do(req *netHttp.Request) ([]byte, error) {
	var (
		resp *netHttp.Response
		err  error
		body []byte
	)
	resp, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var statusError error
	switch resp.StatusCode {
	case netHttp.StatusUnauthorized:
		statusError = &client.ErrorResponse{
			Reason:     "empty or wrong token",
			StatusCode: resp.StatusCode,
		}
	case netHttp.StatusPaymentRequired:
		statusError = &client.ErrorResponse{
			Reason:     "exceeded request limit",
			StatusCode: resp.StatusCode,
		}
	case netHttp.StatusNotFound:
		statusError = &client.ErrorResponse{
			Reason:     "not found",
			StatusCode: resp.StatusCode,
		}
	case netHttp.StatusTooManyRequests:
		statusError = &client.ErrorResponse{
			Reason:     "to many requests",
			StatusCode: resp.StatusCode,
			Retry:      true,
		}
	}
	if statusError != nil {
		c.client.Logger.Error("got client error response", zap.String("err", statusError.Error()))
		return nil, statusError
	}

	body, err = io.ReadAll(resp.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err != nil {
		c.client.Logger.Error("got error on unpacking response", zap.Error(err))
	}

	errorMessage := new(models.ErrorMessage)
	if err = json.Unmarshal(body, errorMessage); err == nil {
		c.client.Logger.Error("got client error message", zap.String("err", errorMessage.Message))
		return nil, fmt.Errorf("%s", errorMessage.Message)
	}

	return body, nil
}

func NewKinopoiskApiUnofficialClient(client *client.Client) *KinopoiskApiUnofficialClient {
	return &KinopoiskApiUnofficialClient{
		client: client,
	}
}
