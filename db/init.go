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

package db

import (
	"encoding/json"
	"log"
	"os"
)

var metadata []MetadataItem

func Init(path string) error {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open metadata.json: %s", err)
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(&metadata)
	if err != nil {
		log.Fatalf("failed to decode metadata.json: %s", err)
	}

	return nil
}
