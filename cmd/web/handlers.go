package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	price := strconv.FormatUint(app.item.GetPrice(), 10)
	w.Header().Add("Content-Type", "text")
	fmt.Fprintf(w, "ITEMS\n-------------\nName: %s\nPrice: %s%s\n", app.item.GetName(), price, app.item.GetCurrency())
}
