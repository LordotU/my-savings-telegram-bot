package fixerio

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type symbols map[string]string
type rates map[string]float64
type ratesNested map[string]rates

type conversionQuery struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
type conversionInfo struct {
	Timestamp int     `json:"timestamp"`
	Rate      float64 `json:"rate"`
}

type FixerIO struct {
	APIKey string
	Base   string
	Secure bool
}

type ResponseSuccess struct {
	Success bool `json:"success"`
}

type ResponseSymbols struct {
	ResponseSuccess
	Symbols symbols `json:"symbols"`
}

type ResponseLatest struct {
	ResponseSuccess
	Timestamp int    `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     rates  `json:"rates"`
}

type ResponseHistorical struct {
	ResponseSuccess
	ResponseLatest
	Historical bool `json:"historical"`
}

type ResponseConversion struct {
	ResponseSuccess
	Query      conversionQuery `json:"query"`
	Info       conversionInfo  `json:"info"`
	Historical bool            `json:"historical"`
	Date       string          `json:"date"`
	Result     float64         `json:"result"`
}

type ResponseTimeseries struct {
	ResponseSuccess
	Timeseries bool        `json:"timeseries"`
	StartDate  string      `json:"start_date"`
	EndDate    string      `json:"end_date"`
	Base       string      `json:"base"`
	Rates      ratesNested `json:"rates"`
}

type ResponseFluctuation struct {
	ResponseSuccess
	Fluctuation bool        `json:"fluctuation"`
	StartDate   string      `json:"start_date"`
	EndDate     string      `json:"end_date"`
	Base        string      `json:"base"`
	Rates       ratesNested `json:"rates"`
}

type ResponseErrorDetail struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

type ResponseError struct {
	ResponseSuccess
	Error ResponseErrorDetail `json:"error"`
}

var BaseUrl = "data.fixer.io"

func GetNew(APIKey string, Base string, Secure bool) (*FixerIO, error) {
	return &FixerIO{
		APIKey,
		Base,
		Secure,
	}, nil
}

func (f *FixerIO) GetSymbols() (*ResponseSymbols, error) {
	url := f.getUrl("symbols", map[string]string{})
	var response ResponseSymbols
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (f *FixerIO) GetLatest(symbols []string) (*ResponseLatest, error) {
	url := f.getUrl("latest", map[string]string{
		"symbols": strings.Join(symbols[:], ","),
	})
	var response ResponseLatest
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (f *FixerIO) GetHistorical(date string, symbols []string) (*ResponseHistorical, error) {
	url := f.getUrl(date, map[string]string{
		"symbols": strings.Join(symbols[:], ","),
	})
	var response ResponseHistorical
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (f *FixerIO) GetConversion(from string, to string, amount float64, date string) (*ResponseConversion, error) {
	url := f.getUrl("convert", map[string]string{
		"from":   from,
		"to":     to,
		"amount": strconv.FormatFloat(float64(amount), 'f', -1, 64),
		"date":   date,
	})
	var response ResponseConversion
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (f *FixerIO) GetTimeseries(start_date string, end_data string) (*ResponseTimeseries, error) {
	url := f.getUrl("timeseries", map[string]string{
		"start_date": start_date,
		"end_data":   end_data,
	})
	var response ResponseTimeseries
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (f *FixerIO) GetFluctuation(start_date string, end_data string) (*ResponseFluctuation, error) {
	url := f.getUrl("fluctuation", map[string]string{
		"start_date": start_date,
		"end_data":   end_data,
	})
	var response ResponseFluctuation
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (f *FixerIO) getUrl(method string, params map[string]string) string {
	var url bytes.Buffer

	if f.Secure {
		url.WriteString("https://")
	} else {
		url.WriteString("http://")
	}

	url.WriteString(BaseUrl)
	url.WriteString("/api")
	url.WriteString("/" + method)
	url.WriteString("?base=")
	url.WriteString(strings.ToUpper(string(f.Base)))
	url.WriteString("&access_key=")
	url.WriteString(string(f.APIKey))

	for name, value := range params {
		if value != "" {
			url.WriteString("&" + name + "=" + value)
		}
	}

	return url.String()
}

func (f *FixerIO) makeRequest(url string, result interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var responseError ResponseError
	err = json.Unmarshal(body, &responseError)
	if err != nil {
		return err
	}

	if !responseError.Success && responseError.Error.Code != 0 {
		return errors.New(responseError.Error.Info)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	return nil
}
