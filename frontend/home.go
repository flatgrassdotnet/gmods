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

package frontend

import (
	"fmt"
	"gmods/db"
	"html/template"
	"net/http"
	"strconv"

	"github.com/xeonx/timeago"
)

type BaseData struct {
	PageType string
	Title    string

	// home
	Items []db.Item

	// search
	Query string
	Tag   string

	// "pagination"
	Shown  int
	Total  int
	Offset int

	// view
	Item db.Item
}

var templateFuncs = template.FuncMap{"sum": func(num ...int) int {
	var i int
	for _, v := range num {
		i += v
	}
	return i
}, "timeago": timeago.English.Format}

func Home(w http.ResponseWriter, r *http.Request) {
	var bd BaseData

	bd.PageType = "home"
	bd.Title = "gmods.org - Garry's Mod Related Files!"

	bd.Query = r.URL.Query().Get("q")
	bd.Tag = r.PathValue("tag")

	if bd.Query == "" && bd.Tag == "" { // no search query or tag, show random
		for range 20 {
			id, err := db.GetRandomItemID()
			if err != nil {
				continue
			}

			d, err := db.GetItem(id)
			if err != nil {
				continue
			}

			bd.Items = append(bd.Items, d)
		}

		bd.Total = db.GetTotal()
	} else { // show search or tag results
		var err error

		if bd.Query != "" { // search qurry
			bd.Items, err = db.GetItemsByName(bd.Query)

		} else { // tag
			bd.Items, err = db.GetItemsByTag(bd.Tag)
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to query downloads: %s", err), http.StatusInternalServerError)
			return
		}

		bd.Total = len(bd.Items)

		// offset
		if r.URL.Query().Get("o") != "" {
			bd.Offset, err = strconv.Atoi(r.URL.Query().Get("o"))
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to decode offset: %s", err), http.StatusInternalServerError)
				return
			}

			if bd.Offset < 0 {
				http.Error(w, "invalid offset", http.StatusInternalServerError)
				return
			}

			start := min(len(bd.Items), bd.Offset)
			end := min(len(bd.Items), start+20)

			bd.Items = bd.Items[start:end]
		}

		// limit to 20
		if len(bd.Items) > 20 {
			bd.Items = bd.Items[:20]
		}
	}

	bd.Shown = len(bd.Items)

	err := t.Execute(w, bd)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
