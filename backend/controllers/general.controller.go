package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"go-scraper/utils"
	"go-scraper/validations"
	"net/http"
	"slices"
	"sync"
)

// Envia las notebooks scrapeadas de Mercadolibre, Fravega y Fullh4rd
func GeneralGetNotebooks(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")   // Se recibe el sort por query params ("asc", "desc", "")
	limit := r.URL.Query().Get("limit") // Se recibe el limite por query params

	scrapSettings := utils.Settings{
		MinRam:     r.URL.Query().Get("minRam"),
		MaxRam:     r.URL.Query().Get("maxRam"),
		MaxInches:  r.URL.Query().Get("maxInches"),
		MinInches:  r.URL.Query().Get("minInches"),
		MinStorage: r.URL.Query().Get("minStorage"),
		MaxStorage: r.URL.Query().Get("maxStorage"),
		Processor:  r.URL.Query().Get("processor"),
		MinPrice:   r.URL.Query().Get("minPrice"),
		MaxPrice:   r.URL.Query().Get("maxPrice"),
	}

	// Se hacen las validaciones (si no se cumplen se envia un error)
	if !validations.ValidateSettings(scrapSettings, w) {
		return
	}

	if !validations.ValidateSort(sort, w) {
		return
	}

	if !validations.ValidateLimit(limit, w) {
		return
	}

	// Si llego hasta aca, ya se valido todo correctamente
	limitNum := utils.GetCorrectLimit(limit)

	// Se usan canales para guardar los resultados al trabajar de forma concurrente
	fullH4rdCh := make(chan []utils.Product)
	mercadolibreCh := make(chan []utils.Product)
	fravegaCh := make(chan []utils.Product)

	// Se usa un WaitGroup para manejar las goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	// Se scrapean las notebooks de los 3 sitios (de forma concurrente, con goroutines)
	go func() {
		defer wg.Done()
		visitUrl := "https://www.fullh4rd.com.ar/cat/search/notebook"
		fullH4rdProducts := scraper.ScrapFullH4rd(visitUrl, scrapSettings)
		fullH4rdCh <- fullH4rdProducts
	}()

	go func() {
		defer wg.Done()
		visitUrl := "https://listado.mercadolibre.com.ar/computacion/laptops-accesorios/notebooks"
		mercadolibreProducts := scraper.ScrapMercadoLibre(visitUrl, scrapSettings)
		mercadolibreCh <- mercadolibreProducts
	}()

	go func() {
		defer wg.Done()
		visitUrl := "https://www.fravega.com/l/informatica/notebooks/?"
		fravegaProducts := scraper.ScrapFravega(visitUrl, scrapSettings)
		fravegaCh <- fravegaProducts
	}()

	// Se cierran los canales una vez que las goroutines hayan terminado
	go func() {
		wg.Wait()
		close(fullH4rdCh)
		close(mercadolibreCh)
		close(fravegaCh)
	}()

	// Traemos los productos de los canales
	fullH4rdProducts := <-fullH4rdCh
	mercadolibreProducts := <-mercadolibreCh
	fravegaProducts := <-fravegaCh

	// Se concatenan los resultados de los productos
	allProducts := append(fullH4rdProducts, mercadolibreProducts...)
	allProducts = append(allProducts, fravegaProducts...)

	// Se ordenan los productos
	if sort == "asc" {
		slices.SortFunc(allProducts, utils.CmpProductAsc)
	} else if sort == "desc" {
		slices.SortFunc(allProducts, utils.CmpProductDesc)
	}

	// Se traen los productos hasta un limite
	allProducts = utils.LimitProducts(limitNum, allProducts)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(allProducts)
}
