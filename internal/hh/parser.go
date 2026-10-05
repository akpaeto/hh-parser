package hh

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var cityAreas = map[string]string{
	"Москва":          "1",
	"Санкт-Петербург": "2",
	"Барнаул":         "70",
}

//http://localhost:8080/

func SearchVacancies(query string, city string) ([]Vacancy, error) {
	searchURL := "https://hh.ru/search/vacancy?text=NAME%3A" + url.QueryEscape(query)

	if area, ok := cityAreas[city]; ok {
		searchURL += "&area" + area
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(searchURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HH вернул статус: %s", resp.Status)

	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	vacancies := make([]Vacancy, 0)

	doc.Find("[data-qa='serp-item__title-text']").Each(func(i int, s *goquery.Selection) {

		card := s.ParentsFiltered("div.vacancy-card--n77Dj8TY8VIUF0yM").First()

		title := strings.TrimSpace(

			card.Find("[data-qa='serp-item__title-text']").Text(),
		)

		company := strings.TrimSpace(
			card.Find("[data-qa='vacancy-serp__vacancy-employer-text']").Text(),
		)

		salary := extractSalary(card)

		location := strings.TrimSpace(

			card.Find("[data-qa='vacancy-serp__vacancy-address']").Text(),
		)

		url, _ := card.Find("a[data-qa='serp-item__title']").Attr("href")

		vacancy := Vacancy{
			Title:    title,
			Company:  company,
			Salary:   salary,
			Location: location,
			URL:      url,
		}

		vacancies = append(vacancies, vacancy)
	})

	return vacancies, nil

}

func extractSalary(card *goquery.Selection) string {

	var salaryParts []string
	card.Find("data").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")
		if !exists {
			return
		}

		number, err := strconv.ParseInt(value, 10, 64)
		if err != nil || number < 1000 {
			return
		}

		text := strings.TrimSpace(s.Text())

		if text != "" {
			salaryParts = append(salaryParts, text)
		}

	})

	if len(salaryParts) == 0 {
		return ""
	}

	return strings.Join(salaryParts, " - ") + "₽"

}
