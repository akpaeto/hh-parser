package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://hh.ru/vacancies/go-razrabotchik"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}

	defer resp.Body.Close()
	fmt.Println("статус:", resp.Status)

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("Ошибка чтения HTML:", err)
		return
	}

	fmt.Println("Ищем PHP-разработчика:")

	fmt.Println("Вакансии:")

	doc.Find("[data-qa='serp-item__title-text']").Each(func(i int, s *goquery.Selection) {
		fmt.Println(i+1, s.Text())
	})

}
