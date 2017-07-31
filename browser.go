package wiphonego

import (
	"fmt"
	"net/http/cookiejar"
	"net/http"

	"net/url"
	"bytes"
)

import (
	"golang.org/x/net/publicsuffix"
	"os"
	"encoding/json"
	"io/ioutil"
	"github.com/go-redis/redis"
	"time"
)


type WebFetcher struct {
	Client *http.Client
	BaseUrl *url.URL

}

func (wb *WebFetcher) SaveCookiesFile(path string) error {
	b, err := json.Marshal(wb.Cookies())
	ioutil.WriteFile(path, b, os.FileMode(0777))
	return err
}

func (wb *WebFetcher) LoadCookiesFile(path string) error {
	a, err := ioutil.ReadFile("cookies.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	var cookies  []*http.Cookie
	json.Unmarshal(a, &cookies)
	wb.Client.Jar.SetCookies(wb.BaseUrl, cookies)
	return err
}

func (wb *WebFetcher) SaveCookiesRedis(key string, client *redis.Client) error {
	fmt.Println(wb.Cookies())
	b, err := json.Marshal(wb.Cookies())
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	s := client.Set(key, string(b), 3600 * time.Second)
	fmt.Println(s.Err())
	return s.Err()
}

func (wb *WebFetcher) LoadCookiesRedis(key string, client *redis.Client) error {
	s := client.Get(key)
	jsonCookies, err := s.Result()
	if err != nil {
		return err
	}

	var cookies  []*http.Cookie
	json.Unmarshal([]byte(jsonCookies), &cookies)
	wb.Client.Jar.SetCookies(wb.BaseUrl, cookies)
	return err
}


func NewWebFetcher(url *url.URL) *WebFetcher {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	cookieJar, _ := cookiejar.New(&options)

	client := &http.Client{
		Jar: cookieJar,
	}

	return &WebFetcher{Client: client, BaseUrl: url}
}

func (wb *WebFetcher) Body(response *http.Response) (string, error) {
	defer response.Body.Close()
	r, err := ioutil.ReadAll(response.Body)
	return string(r), err
}

func (wb *WebFetcher) Get(url string) (*http.Response, error) {
	//"https://yosoymas.masmovil.es/validate/"
	 return wb.Client.Get(url)
}

func (wb *WebFetcher) Cookies() ([]*http.Cookie){
	//r.Request.URL
	return wb.Client.Jar.Cookies(wb.BaseUrl)
}


func (wb *WebFetcher) Post(url string, values url.Values) (*http.Response, error) {
	body := bytes.NewBufferString(values.Encode())
//"https://yosoymas.masmovil.es/validate/"
	rsp, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	//rsp, err := http.NewRequest("POST", "https://httpbin.org/post", body)
	rsp.Header.Set("User-Agent", "Mozilla")
	rsp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//rsp.Header.Set("Cookie", "sid=141295f502e07f041d21801eff4b9384; visid_incap_967703=KFIHTnZTTECZ+q3CZ0Sf7fgeXlkAAAAAQUIPAAAAAAAdLZPA/2B9PFQqHXoL6hnS; incap_ses_504_967703=0A50WubNsgTYi+7YupH+BvgeXlkAAAAA4tmZh7UuOlx/uPF0v2Hz2w==")
	rsp.Header.Set("Accept", "*/*")
	return wb.Client.Do(rsp)
}
