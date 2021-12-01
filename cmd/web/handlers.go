package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	price := strconv.FormatUint(app.item.GetPrice(), 10)
	w.Header().Add("Content-Type", "text")
	fmt.Fprintf(w, "ITEMS\n-------------\nName: %s\nPrice: %s%s\n", app.item.GetName(), price, app.item.GetCurrency())
}

func (app *application) updatePrice(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseUint(string(reqBody), 10, 64)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.item.UpdatePrice(uint64(price))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		app.errorLog.Println(err.Error())
		return
	}
}
