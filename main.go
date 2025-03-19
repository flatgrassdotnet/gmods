/*
	gmods - a rewrite of garrysmod.org
	Copyright (C) 2025  Pancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"gmods/db"
	"gmods/frontend"
	"log"
	"net/http"
)

func main() {
	// flag stuff
	port := flag.Int("port", 80, "http listen port")
	flag.Parse()

	// set up frontend
	err := frontend.Init()
	if err != nil {
		log.Fatalf("failed to initialize frontend: %s", err)
	}

	// set up database
	err = db.Init("metadata.json")
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err)
	}

	// set http routes
	http.HandleFunc("GET /", frontend.Home)
	http.HandleFunc("GET /view/{id}", frontend.View)

	http.HandleFunc("POST /download/{id}", frontend.Download)

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// start http server
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
