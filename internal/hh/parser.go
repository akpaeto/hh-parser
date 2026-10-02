package hh

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SearchVacancies(query string) ([]Vacancy, error) {
	url := "https://hh.ru/vacancies/" + strings.ReplaceAll(query, " ", "-")

	resp, err := http.Get(url)

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

	var vacancies []Vacancy

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

		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
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
