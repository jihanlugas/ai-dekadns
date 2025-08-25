package helper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type Myjar struct {
	jar map[string][]*http.Cookie
}

func (p *Myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	// fmt.Printf("The URL is : %s\n", u.String())
	// fmt.Printf("The cookie being set is : %s\n", cookies)
	p.jar[u.Host] = cookies
}

func (p *Myjar) Cookies(u *url.URL) []*http.Cookie {
	// fmt.Printf("The URL is : %s\n", u.String())
	// fmt.Printf("Cookie being returned is : %s\n", p.jar[u.Host])
	return p.jar[u.Host]
}

var client = &http.Client{}

func PostWithHeader(url string, body []byte, header map[string]string) ([]byte, error) {
	return doRequestWithHeader("POST", url, body, header)
}

func GetWithQuery(url string, queries map[string]string, header map[string]string) ([]byte, error) {
	return doRequestWithQuery("GET", url, queries, header)
}

func GetWithHeader(url string, body []byte, header map[string]string) ([]byte, error) {
	return doRequestWithHeader("GET", url, body, header)
}

func doRequestWithQuery(method string, uri string, queries map[string]string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, uri, nil)
	queryData := url.Values{}

	// log.Infof("LOG REQ: %v", req)
	for k, v := range header {
		req.Header.Set(k, v)
	}

	for k, v := range queries {
		if v != "" {
			queryData.Add(k, v)
		}
	}

	req.URL.RawQuery = queryData.Encode()

	// log.Infof("LOG ERR: %v", err)
	if err != nil {
		return nil, err
	}

	log.Infof("LOG REQ: %v", req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// log.Infof("HIT API RESULT:")
	// log.Infof("URL: %v", uri)
	log.Infof("STATUS CODE: %v", resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	if isStatusError(resp.StatusCode) {
		return nil, fmt.Errorf("status error: %v %v", resp.StatusCode, string(respBody))
	}

	defer resp.Body.Close()
	return respBody, err
}

func doRequestWithHeader(method string, url string, body []byte, header map[string]string) ([]byte, error) {
	payload := strings.NewReader(string(body))
	req, err := http.NewRequest(method, url, payload)
	log.Infof("LOG REQ: %v", req)
	for k, v := range header {
		req.Header.Set(k, v)
	}

	log.Infof("LOG ERR: %v", err)
	log.Println("body request ", payload)

	if err != nil {
		return nil, err
	}

	log.Infof("LOG REQ: %v", req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Infof("HIT API RESULT:")
	log.Infof("URL: %v", url)
	log.Infof("STATUS CODE: %v", resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	if isStatusError(resp.StatusCode) {
		return nil, fmt.Errorf("status error: %v %v", resp.StatusCode, string(respBody))
	}

	defer resp.Body.Close()
	return respBody, err
}

func isStatusError(statusCode int) bool {
	return statusCode >= http.StatusBadRequest
}

type paramsState struct {
	params map[string]interface{}

	err error
}

func (s *paramsState) Allowed(allowed []string) *paramsState {
	allowed = append(allowed, "page")
	allowed = append(allowed, "limit")
	err := AllowedKey(s.params, allowed)
	if err != nil {
		s.err = errors.New("parameters with '" + err.Error() + "' is not allowed")
		return s
	}
	return s
}

func (s *paramsState) AllowedAndCast(key string, dataType reflect.Kind) *paramsState {
	s.Allowed([]string{key})
	if s.err != nil {
		return s
	}
	switch dataType {
	case reflect.Int:
		if s.params[key] != nil {
			data, ok := s.params[key].(string)
			if !ok {
				s.err = errors.New(fmt.Sprintf("parameter %s cant cast to ", key))
				return s
			}
			res, err := strconv.Atoi(data)
			if err != nil {
				s.err = errors.New(fmt.Sprintf("parameter %s cant cast to int", key))
				return s
			}
			s.params[key] = res
		}
	case reflect.Int64:
		if s.params[key] != nil {
			data, ok := s.params[key].(string)
			if !ok {
				s.err = errors.New(fmt.Sprintf("parameter %s cant cast to int64", key))
				return s
			}
			res, err := strconv.Atoi(data)
			if err != nil {
				s.err = errors.New(fmt.Sprintf("parameter %s cant cast to int64", key))
				return s
			}
			s.params[key] = int64(res)
		}
	case reflect.Float64:
		if s.params[key] != nil {
			data, ok := s.params[key].(string)
			if !ok {
				s.err = errors.New(fmt.Sprintf("parameter %s cant cast to float64", key))
				return s
			}
			res, err := strconv.Atoi(data)
			if err != nil {
				s.err = errors.New(fmt.Sprintf("parameter %s cant cast to float64", key))
				return s
			}
			s.params[key] = float64(res)
		}
	default:
		s.err = errors.New("not support parameter check data type")
	}
	return s
}

func (s *paramsState) Must(mustList []string) *paramsState {
	s.Allowed(mustList)
	if s.err != nil {
		return s
	}
	for _, item := range mustList {
		var exist bool
		for key, _ := range s.params {
			if key == item {
				exist = true
				break
			}
		}
		if !exist {
			s.err = errors.New("parameters with '" + item + "' is needed")
		}
	}
	return s
}

func (s *paramsState) BlackList(blocked []string) *paramsState {
	for field, _ := range s.params {
		for _, allowField := range blocked {
			if field == allowField {
				s.err = errors.New("parameters with '" + field + "' is blocked")
				return s
			}
		}
	}
	return s
}

func (s *paramsState) Result() (map[string]interface{}, error) {
	return s.params, s.err
}

func ParamsToMapOneValue(c *gin.Context) *paramsState {
	params := map[string]interface{}{}
	for field, item := range c.Request.URL.Query() {
		var value string
		for _, val := range item {
			value = val
			break
		}
		params[field] = value
	}
	return &paramsState{params: params}
}
