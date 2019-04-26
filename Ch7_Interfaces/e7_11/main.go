/*
Add additional handlers so that clients can create, read, update, and delete database entries.
For example, a request of the form /update?item=socks&price=6 will update the price of an item in the
inventory and report an error if the item does not exist or if the price is invalid. (Warning: this
change introduces concurrent variable updates.)
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	addr := "localhost:8888"

	db := database{
		items: map[string]pounds{"shoes": 50, "socks": 5},
	}

	http.HandleFunc("/", db.index)
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/read", db.read)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/delete", db.delete)

	log.Println("Server running at: ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type pounds float32

func (p pounds) String() string {
	return fmt.Sprintf("£%.2f", p)
}

type database struct {
	sync.Mutex
	items map[string]pounds
}

func (db *database) index(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/list", http.StatusPermanentRedirect)
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()

	for item, price := range db.items {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.Lock()
	defer db.Unlock()

	if price, ok := db.items[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.Lock()
	defer db.Unlock()

	if _, ok := db.items[item]; ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "item exists: %s\n", item)
		return
	}
	p := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(p, 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid price: %s\n", p)
		return
	}
	db.items[item] = pounds(price)
	fmt.Fprintf(w, "%s: %s\n", item, db.items[item])
}

func (db *database) read(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.Lock()
	defer db.Unlock()

	if _, ok := db.items[item]; ok {
		fmt.Fprintf(w, "%s: %s\n", item, db.items[item])
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.Lock()
	defer db.Unlock()

	p := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(p, 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid price: %s\n", p)
		return
	}

	if _, ok := db.items[item]; ok {
		db.items[item] = pounds(price)
		fmt.Fprintf(w, "%s\n", db.items[item])
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.Lock()
	defer db.Unlock()

	if _, ok := db.items[item]; ok {
		delete(db.items, item)
		if _, ok := db.items[item]; !ok {
			fmt.Fprintf(w, "deleted item: %q\n", item)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
