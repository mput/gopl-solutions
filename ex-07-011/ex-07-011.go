package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

type dollars int64

func parseDollars(s string) (dollars, error) {
	r := regexp.MustCompile("^(\\d+)(?:\\.(\\d{1,2}))?$")
	res := r.FindStringSubmatch(s)
	if len(res) < 2 {
		return 0, fmt.Errorf("parse dollar error: %s is wrong value", s)
	}
	whole, err := strconv.Atoi(res[1])
	if err != nil {
		panic(err)
	}
	var partial int = 0
	if len(res[2]) == 1 {
		partial, err = strconv.Atoi(res[2])
		partial *= 10
	}
	if len(res[2]) == 2 {
		partial, err = strconv.Atoi(res[2])
	}
	return dollars((whole * 100) + partial), nil
}

func (d dollars) String() string {
	t := fmt.Sprintf("%03d", d)
	return fmt.Sprintf("%s.%s", t[:len(t)-2], t[len(t)-2:])
}

type database struct {
	Data map[string]dollars
	Lock sync.Mutex
}

func (db *database) Price(name string) dollars {
	return db.Data[name]
}

func (db *database) Prices() map[string]dollars {
	return db.Data
}

func (db *database) SetPrice(name string, price dollars) {
	db.Lock.Lock()
	defer db.Lock.Unlock()
	db.Data[name] = price
}

func (db *database) RenderTable(w http.ResponseWriter) {
	var t = template.Must(template.ParseFiles("table.html"))
	err := t.Execute(w, struct {
		Data map[string]dollars
	}{Data: db.Data})
	if err != nil {
		log.Fatal(err)
	}
}

func (db *database) List(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		db.RenderTable(w)
	case "POST":
		name := r.FormValue("name")
		priceStr := r.FormValue("price")
		price, err := parseDollars(priceStr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if name == "" {
			http.Error(w, "Name is required", 400)
			return
		}
		db.SetPrice(name, price)
		db.RenderTable(w)
	}

}

func main() {
	data := map[string]dollars{"Crocs": dollars(1999)}
	datadb := database{Data: data}

	http.HandleFunc("/list", datadb.List)

	log.Fatal(http.ListenAndServe("localhost:9999", nil))
}
