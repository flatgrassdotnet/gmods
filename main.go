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
	"gmods/db"
	"gmods/frontend"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	// flag stuff
	dbuser := flag.String("dbuser", "gmods", "database user's name")
	dbpass := flag.String("dbpass", "", "database user's password")
	dbproto := flag.String("dbproto", "tcp", "database connection protocol")
	dbaddr := flag.String("dbaddr", "localhost", "database server address")
	dbname := flag.String("dbname", "gmods", "database name")
	proto := flag.String("proto", "tcp", "proto for web server")
	addr := flag.String("addr", "127.0.0.1:80", "address for web server")
	flag.Parse()

	// set up frontend
	err := frontend.Init()
	if err != nil {
		log.Fatalf("failed to initialize frontend: %s", err)
	}

	// set up database
	err = db.Init(*dbuser, *dbpass, *dbproto, *dbaddr, *dbname)
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err)
	}

	// set http routes
	http.HandleFunc("GET /", frontend.Home)
	http.HandleFunc("GET /tag/{tag}", frontend.Home)
	http.HandleFunc("GET /view/{id}", frontend.View)

	http.HandleFunc("POST /download/{id}", frontend.Download)

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// start http server
	if *proto == "unix" {
		err = os.Remove(*addr)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("failed to delete unix socket: %s", err)
		}
	}

	l, err := net.Listen(*proto, *addr)
	if err != nil {
		log.Fatalf("failed to create web server listener: %s", err)
	}

	defer l.Close()

	if *proto == "unix" {
		err = os.Chmod(*addr, 0777)
		if err != nil {
			log.Fatalf("failed to set unix socket permissions: %s", err)
		}
	}

	http.Serve(l, nil)
}
