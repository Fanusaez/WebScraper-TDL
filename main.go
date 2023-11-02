package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Url   string `json:"url"`
}

func main() {
	// Se solicita al usuario que producto desea buscar
	var productToSearch string

	fmt.Print("Ingresa un producto a buscar: ")
	fmt.Scanln(&productToSearch)

	var visitUrl string = "https://listado.mercadolibre.com.ar/" + productToSearch

	// Crea una nueva instancia de Colly
	c := colly.NewCollector()

	// Slice para almacenar los productos scrapeados
	var products []Product

	c.OnHTML(".ui-search-result__wrapper", func(e *colly.HTMLElement) {
		product := Product{
			Name:  e.ChildText(".ui-search-item__title"),
			Price: e.ChildText("div.ui-search-price__second-line span.andes-money-amount__fraction"),
			Url:   e.ChildAttr("a", "href"),
		}

		products = append(products, product)
	})

	// Se visita el sitio y se imprimen los productos scrapeados
	c.Visit(visitUrl)

	for _, product := range products {
		fmt.Println("Nombre:", product.Name)
		fmt.Println("Precio:", product.Price)
		fmt.Println("Url:", product.Url)
	}

	saveAsJsonFile(products) // Guardamos los datos en un archivo
}

// Guarda un slice de productos como json
func saveAsJsonFile(products []Product) {
	data, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("products.json", data, 0644)
}
