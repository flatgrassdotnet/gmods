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
)

func View(w http.ResponseWriter, r *http.Request) {
	// TODO: do this once and store the template
	t, err := template.New("base.html").Funcs(templateFuncs).ParseFiles("templates/base.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create template: %s", err), http.StatusInternalServerError)
		return
	}

	t, err = t.ParseGlob("templates/include/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create template: %s", err), http.StatusInternalServerError)
		return
	}

	var bd BaseData

	bd.PageType = "view"

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode id value: %s", err), http.StatusInternalServerError)
		return
	}

	bd.Item, err = db.GetItem(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get download: %s", err), http.StatusInternalServerError)
		return
	}

	bd.Title = fmt.Sprintf("%s - Download!", bd.Item.Name)

	err = t.Execute(w, bd)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
