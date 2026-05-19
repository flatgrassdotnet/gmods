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
	"embed"
	"io/fs"
	"text/template"
)

var t *template.Template

var (
	//go:embed templates
	templates      embed.FS
	TemplatesFS, _ = fs.Sub(templates, "templates")

	//go:embed assets
	assets      embed.FS
	AssetsFS, _ = fs.Sub(assets, "assets")
)

func Init() error {
	var err error

	t, err = template.New("base.html").Funcs(templateFuncs).ParseFS(TemplatesFS, "base.html")
	if err != nil {
		return err
	}

	t, err = t.ParseFS(TemplatesFS, "include/*.html")
	if err != nil {
		return err
	}

	return nil
}
