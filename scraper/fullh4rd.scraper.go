package scraper

import (
	"fmt"
	"go-scraper/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Scrapea datos de fullh4rd, a partir de una url, y devuelve los productos
func ScrapFullH4rd(url string, scrapSettings utils.Settings) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product
	// Se aplican los settings/filtros de scrapeo
	visitUrl := applyScrapSettingsFullH4rd(url, &scrapSettings)

	minPrice, err := strconv.ParseFloat(scrapSettings.MinPrice, 64)
	if err != nil {
		minPrice = 0

	}

	maxPrice, err := strconv.ParseFloat(scrapSettings.MaxPrice, 64)
	if err != nil {
		maxPrice = 0
	}

	// Se define el comportamiento al scrapear
	c.OnHTML("div[class='item product-list']", func(e *colly.HTMLElement) {
		price := e.ChildText("div.price")
		elements := strings.Fields(price)
		if len(elements) >= 2 {
			price = elements[0]
		}

		productName := e.ChildText("h3")
		specs := parseSpecs(productName)

		product := utils.Product{
			Name:   productName,
			Price:  utils.ConvertPriceToNumber(price),
			Url:    "https://www.fullh4rd.com.ar/" + e.ChildAttr("a", "href"),
			Origin: "FullH4rd",
			Specs:  specs,
		}

		if verifyProductFullH4rd(product.Name, &scrapSettings) && len(specs.Processor) > 7 {

			if maxPrice > 0 || minPrice > 0 {
				if maxPrice != 0 && minPrice == 0 {
					if product.Price < int(maxPrice) {
						products = append(products, product)
					}
				} else if maxPrice == 0 && minPrice != 0 {
					if product.Price > int(minPrice) {
						products = append(products, product)
					}
				} else if maxPrice != 0 && minPrice != 0 {
					if product.Price > int(minPrice) && product.Price < int(maxPrice) {
						products = append(products, product)
					}
				}

			} else {
				products = append(products, product)
			}

		}

	})

	c.Visit(visitUrl) // Se visita el sitio a scrapear

	return products
}

// Funcion auxiliar, aplica los settings de busqueda sobre una url para fullh4rd y devuelve
// una nueva url
func applyScrapSettingsFullH4rd(url string, scrapSettings *utils.Settings) string {
	appendStr := "%20"

	// Se aplican los settings para scrapear
	if scrapSettings.Ram != "" {
		url += fmt.Sprintf("%s%sgb", appendStr, scrapSettings.Ram)
	}

	if scrapSettings.Storage != "" {

		if scrapSettings.Storage == "1000" {

			url += fmt.Sprintf("%s1tb", appendStr)
		} else {

			url += fmt.Sprintf("%s%sgb", appendStr, scrapSettings.Storage)
		}

	}

	if scrapSettings.Inches != "" {

		floatValue, err := strconv.ParseFloat(scrapSettings.Inches, 64)
		if err != nil {
			fmt.Println("Error parsing float:", err)
		}
		intValue := int(floatValue)
		result := strconv.Itoa(intValue)

		url += fmt.Sprintf("%s%s", appendStr, result)
	}

	if scrapSettings.Processor != "" {

		if scrapSettings.Processor == "amd" {
			url += fmt.Sprintf("%s%s", appendStr, "ryzen")
		} else if scrapSettings.Processor == "apple" {
			url += fmt.Sprintf("%s%s", appendStr, "apple")
		}

	}

	if scrapSettings.MinPrice != "" || scrapSettings.MaxPrice != "" {
		if scrapSettings.MinPrice == "" {
			scrapSettings.MinPrice = "0"
		}

		if scrapSettings.MaxPrice == "" {
			scrapSettings.MaxPrice = "0"
		}
	}

	// Se crea y devuelve la url que finalmente se va a scrapear
	visitUrl := url

	return visitUrl
}

// Funcion auxiliar, realiza validaciones para el scrapeo de productos de fullh4rd
func verifyProductFullH4rd(name string, scrapSettings *utils.Settings) bool {

	lowerName := strings.ToLower(name)

	if (strings.Contains(lowerName, "notebook") || strings.Contains(lowerName, "laptop")) &&
		!(strings.Contains(lowerName, "cooler")) {

		// intel, amd o apple
		if !(scrapSettings.Processor == "") {

			if scrapSettings.Processor == "amd" {

				return strings.Contains(lowerName, "ryzen")
			} else if scrapSettings.Processor == "apple" {

				return strings.Contains(lowerName, "apple")
			} else if scrapSettings.Processor == "intel" {

				return !strings.Contains(lowerName, "apple") && !strings.Contains(lowerName, "amd") && !strings.Contains(lowerName, "ryzen")
			}

		}

		if !(scrapSettings.Ram == "") {
			ram := scrapSettings.Ram + `gb`
			if !strings.Contains(lowerName, ram) {
				return false
			}
		}

		if !(scrapSettings.Storage == "") {

			storage := scrapSettings.Storage + `gb`
			if scrapSettings.Storage == "1000" {
				storage = "1tb"
			}

			if !strings.Contains(lowerName, storage) {
				return false
			}
		}

		if !(scrapSettings.Inches == "") {
			inches := scrapSettings.Inches + `"`
			if !strings.Contains(lowerName, inches) {
				return false
			}
		}

		return true

	} else {
		return false
	}
}

func parseSpecs(input string) utils.Specs {
	var specs utils.Specs

	ramRegex := regexp.MustCompile(`(\d+)GB`)
	storageRegex := regexp.MustCompile(`(\d+)(GB|TB|G)`)

	ramMatches := ramRegex.FindAllStringSubmatch(input, -1)
	storageMatches := storageRegex.FindAllStringSubmatch(input, -1)

	maxRam := 0
	for _, match := range ramMatches {
		ram, err := strconv.Atoi(match[1])
		if err == nil && ram > maxRam {
			maxRam = ram
		}
	}

	for _, match := range ramMatches {
		ram, _ := strconv.Atoi(match[1])
		if ram == maxRam {
			specs.Ram = match[0]
		}
	}

	for _, match := range storageMatches {
		if specs.Ram == "" || match[0] != specs.Ram {
			specs.Storage = match[0]
		}
	}

	if !strings.Contains(specs.Storage, "TB") {

		specs.Ram, specs.Storage = specs.Storage, specs.Ram
	}

	if strings.HasSuffix(specs.Ram, "G") {

		specs.Ram, specs.Storage = specs.Storage, specs.Ram
	}

	displayPattern := regexp.MustCompile(`(\d+(\.\d+)?)\"`)

	match := displayPattern.FindStringSubmatch(input)

	if len(match) > 1 {
		specs.Inches = match[1]
	}

	if strings.Contains(input, "RYZEN") {

		substrings := strings.Fields(input)

		result := "RYZEN"

		include := false

		for _, substring := range substrings {

			if strings.Contains(substring, "GB") {
				break
			}

			if include {
				result += " " + substring
			}

			if strings.Contains(substring, "RYZEN") {
				include = true
			}
		}

		result = strings.TrimSpace(result)

		specs.Processor = result

	} else {
		result := ""
		re := regexp.MustCompile(`(?:I[0-9]+-[0-9A-Za-z]+)|(?:I[0-9]+\s[0-9A-Za-z]+)`)

		match := re.FindString(input)

		result = match
		specs.Processor = result
	}

	if strings.Contains(specs.Processor, "RTX") {
		specs.Processor = removeGPUFromProcessor(specs.Processor)
	}

	return specs
}

func removeGPUFromProcessor(processor string) string {

	r := regexp.MustCompile(`RTX\s?\d+`)

	result := r.ReplaceAllString(processor, "")

	trimmed := strings.TrimSpace(result)

	if len(trimmed) < len(result) {

		trimmed = trimmed[:len(trimmed)-1]
	}

	return trimmed
}
