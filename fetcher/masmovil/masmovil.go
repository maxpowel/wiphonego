package masmovil

import (
	"net/url"
	"time"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strconv"
	"github.com/maxpowel/wiphonego"
	"github.com/maxpowel/dislet"
	"strings"
	"github.com/go-redis/redis"
	"fmt"
	"math"
)

type Fetcher struct {
	Fetcher *wiphonego.WebFetcher
	Credentials *wiphonego.Credentials
	Kernel *dislet.Kernel
}

func (f *Fetcher) isLogged() (error) {
	time.Sleep(1 * time.Second)
	r, err := f.Fetcher.Get("https://yosoymas.masmovil.es/")
	if err != nil {
		return err
	}
	body, err := f.Fetcher.Body(r)
	if err != nil {
		return err
	}

	if strings.Index(body, "Desconectar") > 0 {
		return nil
	} else {
		return fmt.Errorf("Invalid credentials")
	}

}

func (f *Fetcher) login() (error) {

	f.Fetcher.Get("https://yosoymas.masmovil.es/validate/")

	time.Sleep(3 * time.Second)
	form := url.Values{}
	form.Add("action", "login")
	form.Add("url", "")
	form.Add("user", f.Credentials.Username)
	form.Add("password", f.Credentials.Password)
	f.Fetcher.Post("https://yosoymas.masmovil.es/validate/", form)
	time.Sleep(1 * time.Second)


	err := f.isLogged()
	if err != nil {
		return err
	} else {
		client := f.Kernel.Container.MustGet("redis").(*redis.Client)
		f.Fetcher.SaveCookiesRedis(f.CredentialsKey(), client)
		return nil
	}

}

func (f *Fetcher) GetInternetConsumption(phoneNumber string) (wiphonego.UserDeviceConsumption, error){
	err := f.isLogged()
	if err != nil {
		fmt.Println("Loging in")
		loginError := f.login()
		if loginError != nil {
			return wiphonego.UserDeviceConsumption{}, loginError
		}
	} else {
		fmt.Print("Already logged in")
	}

	//f.fetcher.LoadCookies("cookies.json")
	res, err := f.Fetcher.Get("https://yosoymas.masmovil.es/consumo/?line="+phoneNumber)



	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}

	//ci := make(chan InternetConsumption)
	c := wiphonego.UserDeviceConsumption{}
	doc.Find(".box-main-content").Find(".progress").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		//c := s.Find("span").Text()
		re := regexp.MustCompile("([0-9]+)|(infinito)")
		r := re.FindAllString(s.Text(), -1)

			//fmt.Printf("Megas gastados %v de %v\n", r[0], r[1])
			consumed, err := strconv.ParseInt(r[0],10, 64)
			if err == nil {
				if i == 0 {
					c.InternetConsumed = consumed * 1024 * 1024
				} else {
					c.CallConsumed = int(consumed)
				}
			}

			total, err := strconv.ParseInt(r[1],10, 64)
			if err != nil {
				total = math.MaxInt32
			}

			if i == 0 {
				c.InternetTotal = total * 1024 * 1024
			} else {
				c.CallTotal = int(total)
			}

			fmt.Printf("Gastados %v de %v\n", consumed, total)
	})

	//c := <- ci
	return c, nil
	//f.fetcher.LoadCookies("cookies.json")

}

func (f *Fetcher) CredentialsKey() (string){
	return f.Credentials.Username + ":" + f.Credentials.Operator.Name
}

func NewFetcher (credentials *wiphonego.Credentials, kernel *dislet.Kernel) *Fetcher{
	f := &Fetcher{
		Fetcher: wiphonego.NewWebFetcher(&url.URL{Host:"yosoymas.masmovil.es", Scheme:"https"}),
		Credentials: credentials,
		Kernel: kernel,
	}

	client := kernel.Container.MustGet("redis").(*redis.Client)

	f.Fetcher.LoadCookiesRedis(f.CredentialsKey(), client)
	return f
}