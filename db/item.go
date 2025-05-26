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
	"errors"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

type Item struct {
	ID   int
	Name string

	AuthorName string
	AuthorID   int

	Description string
	Tags        []string

	Images []string

	Size      string
	Posted    time.Time
	Downloads int
}

func GetTotal() int {
	return len(metadata)
}

// returns the id of a random download
func GetRandomItemID() (int, error) {
	id, err := metadata[rand.IntN(len(metadata))].GetID()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// returns a download by its id
func GetItem(id int) (Item, error) {
	for _, m := range metadata {
		mid, err := m.GetID()
		if err != nil {
			continue
		}

		if mid != id {
			continue
		}

		var item Item

		item.ID = mid
		item.Name = m.Name

		item.AuthorName = m.GetUploaderName()
		item.AuthorID, err = m.GetUploaderID()
		if err != nil {
			return Item{}, err
		}

		item.Description = strings.ReplaceAll(m.Description, `\n`, "\n")
		item.Tags = m.Tags

		for _, i := range m.Images {
			item.Images = append(item.Images, strings.ReplaceAll(i, "filecache.garrysmods.org", "data.gmods.org"))
		}

		item.Size, err = m.GetPrettySize()
		if err != nil {
			return Item{}, err
		}

		item.Posted = m.GetUploadTime()

		return item, nil
	}

	return Item{}, errors.New("download id not found")
}

// returns downloads with name that contains string
func GetItemsByName(name string) ([]Item, error) {
	var results []Item
	for _, m := range metadata {
		if !strings.Contains(strings.ToLower(m.Name), strings.ToLower(name)) {
			continue
		}

		id, err := m.GetID()
		if err != nil {
			return nil, err
		}

		dl, err := GetItem(id)
		if err != nil {
			return nil, err
		}

		results = append(results, dl)
	}

	return results, nil
}

// return downloads that contain tag
func GetItemsByTag(tag string) ([]Item, error) {
	var results []Item
	for _, m := range metadata {
		if !slices.Contains(m.Tags, tag) {
			continue
		}

		id, err := m.GetID()
		if err != nil {
			return nil, err
		}

		dl, err := GetItem(id)
		if err != nil {
			return nil, err
		}

		results = append(results, dl)
	}

	return results, nil
}
