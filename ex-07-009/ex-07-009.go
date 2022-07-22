package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type Song struct {
	Title   string
	Artist string
	Year int
}

type SongOrderParam struct {
	Title string
	Query string
	Order int
	Position int
}

type SortableSongs struct {
	Songs  []*Song
	OrderParams []order

}

func (s SortableSongs) Len() int {
	return len(s.Songs)
}

func (s SortableSongs) Swap(i, j int) {
	s.Songs[i], s.Songs[j] = s.Songs[j], s.Songs[i]
}

func (s SortableSongs) Less(i, j int) bool {
	for _, e := range s.OrderParams {
		switch e.field {
		case "title":
			if s.Songs[i].Title != s.Songs[j].Title {
				if e.order == 0 {
					return s.Songs[i].Title < s.Songs[j].Title
				}
				return s.Songs[i].Title > s.Songs[j].Title
			}
		case "artist":
			if s.Songs[i].Artist != s.Songs[j].Artist {
				if e.order == 0 {
					return s.Songs[i].Artist < s.Songs[j].Artist
				}
				return s.Songs[i].Artist > s.Songs[j].Artist
			}
		case "year":
			if s.Songs[i].Year != s.Songs[j].Year {
				if e.order == 0 {
					return s.Songs[i].Year < s.Songs[j].Year
				}
				return s.Songs[i].Year > s.Songs[j].Year
			}
		}

	}
	return false

}


func formatOrders(os []order) string{
	res := make([]string, 0, len(os))
	for _,o := range os {
		s := ""
		s = s + "order=" + o.field
		if o.order > 0 {
			s = s + "[desc]"
		}
		res = append(res, s)
	}
	return "?" + strings.Join(res, "&")
}


func buildSongOrderParam (field string, orders []order) SongOrderParam {
	nextOrders := make([]order, 0, 3)
	curPosition := -1
	for i,e := range orders {
		if e.field == field {
			curPosition = i
		}
	}

	curOrder    := -1
	if curPosition > -1 {
		 curOrder = orders[curPosition].order
	}

	if len(orders) == 0 {
		nextOrders = append(nextOrders, order{field, asc})
	} else if curPosition == -1 {
		nextOrders = append(nextOrders, order{field, asc})
		nextLen := 1;
		if len(orders) + 1 >= 3 {
			nextLen = 3
		} else {
			nextLen = len(orders) + 1
		}
		nextOrders = nextOrders[0:nextLen]
		copy(nextOrders[1:], orders)
	} else if curPosition == 0 {
		nextOrders = nextOrders[:len(orders)]
		copy(nextOrders, orders)
		if nextOrders[0].order == asc {
			nextOrders[0].order = desc
		} else {
			nextOrders[0].order = asc
		}
	} else {
		nextOrders = append(nextOrders, order{field, asc})
		for i, e := range orders {
			if e.field != field && i < 2 {
				nextOrders = append(nextOrders, e)
			}
		}
	}

	return SongOrderParam{
		Title: strings.ToUpper(field),
		Query: formatOrders(nextOrders),
		Position: curPosition,
		Order: curOrder,
	}

}

func buildSongOrderParams (o []order) []SongOrderParam {
	res := make([]SongOrderParam,0, 3)
	for _, param := range []string{"title", "artist", "year"} {
		res = append(res, buildSongOrderParam(param, o))
	}
	return res
}

func renderTable (wr io.Writer,s []*Song, os []order) {
	var t = template.Must(template.ParseFiles("table.html"))
	data := struct{
		Songs []*Song
		Order []SongOrderParam
	}{
		Songs: s,
		Order: buildSongOrderParams(os),
	}
	err := t.Execute(wr, data)
	if err != nil {
		log.Fatal(err)

	}
}

const (
	asc = iota
	desc = iota
)

type order struct {
	field string
	order int
}

func parseSortQuery (s string) (order, error) {
	r := regexp.MustCompile("(\\w+)+(?:\\[(asc|desc)\\])?")
	match := r.FindStringSubmatch(s)
	if len(match) < 2  {
		return order{}, fmt.Errorf("wrong order parameter: %s", s)
	}
	ord := asc
	if match[2] == "desc" {
		ord = desc
	}
	return order{match[1], ord}, nil
}


func serveTable(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	ordersRaw := params["order"]
	orders := make([]order, 0, len(ordersRaw))
	for _, or := range ordersRaw {
		order, err := parseSortQuery(or)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		orders = append(orders, order)
	}

	songs := []*Song{
		{"Do it all again", "Z Loud", 2013},
		{"Cross you mind", "A David Solomon", 2018},
		{"High Hopes", "X Pink Floyd", 2013},
		{"Nobody", "B Catello", 2014},
	}


	sort.Sort(SortableSongs{songs, orders})
	renderTable(w, songs, orders)

}

func main() {
	http.HandleFunc("/", serveTable)
	http.ListenAndServe("localhost:8989", nil)

}
