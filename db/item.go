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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Item struct {
	ID          int
	Name        string
	Filename    string
	Description sql.NullString
	Size        int
	Uploader    sql.NullString
	Uploaded    sql.NullTime

	Tags   []string
	Images map[int]string

	Downloads int
	Views     int
}

var ErrInvalidID = errors.New("download id not found")

func (i Item) PrettySize() string {
	return fmt.Sprintf("%.02f MB", float64(i.Size)/1024/1024)
}

func GetItemList(ctx context.Context, tag string, query string) ([]Item, error) {
	q := `SELECT p.id, p.name, p.filename, p.description, p.size, p.uploader, COALESCE(p.uploaded, f.modified), s.downloads, s.views 
	FROM packages p 
	JOIN stats s ON p.id = s.pid 
	LEFT JOIN (
	SELECT pid, MAX(modified) AS modified 
	FROM files 
	WHERE modified > TIMESTAMP'2006-01-01 00:00:00' 
	AND modified < TIMESTAMP'2015-01-01 00:00:00' 
	GROUP BY pid) f 
	ON p.id = f.pid`
	var args []any

	if tag != "" {
		q += " JOIN tags t ON t.pid = p.id WHERE t.tag = ?"
		args = append(args, tag)
	}
	if query != "" {
		// filenames
		q += ` LEFT JOIN ( 
		SELECT DISTINCT pid 
		FROM files 
		WHERE MATCH(path) AGAINST(?) 
		) f2 
		ON f2.pid = p.id 
		WHERE f2.pid IS NOT NULL`
		args = append(args, query)

		// description
		q += " OR MATCH(p.description) AGAINST(?)"
		args = append(args, query)

		// upload name
		// LIKE is needed for handling underscores
		q += " OR p.name LIKE CONCAT('%', ?, '%')"
		args = append(args, strings.ReplaceAll(query, " ", "_"))
	}
	if tag == "" && query == "" {
		q += " ORDER BY RAND() LIMIT 20"
	}

	var items []Item
	rows, err := conn.QueryContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidID
		}

		return nil, err
	}
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.Name, &item.Filename, &item.Description, &item.Size, &item.Uploader, &item.Uploaded, &item.Downloads, &item.Views)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func GetItem(ctx context.Context, id int) (Item, error) {
	item := Item{ID: id, Images: make(map[int]string)}

	q := `SELECT p.name, p.filename, p.description, p.size, p.uploader, COALESCE(p.uploaded, f.modified) 
	FROM packages p 
	LEFT JOIN (
	SELECT pid, MAX(modified) AS modified 
	FROM files 
	WHERE modified > TIMESTAMP'2006-01-01 00:00:00' 
	AND modified < TIMESTAMP'2015-01-01 00:00:00' 
	GROUP BY pid) f 
	ON p.id = f.pid 
	WHERE p.id = ?`

	err := conn.QueryRowContext(ctx, q, id).Scan(&item.Name, &item.Filename, &item.Description, &item.Size, &item.Uploader, &item.Uploaded)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return item, ErrInvalidID
		}

		return item, err
	}

	rows, err := conn.QueryContext(ctx, "SELECT i.id, i.res FROM images i JOIN packages p ON p.id = i.pid WHERE p.id = ?", id)
	if err != nil {
		return item, err
	}
	for rows.Next() {
		var id int
		var res string
		err = rows.Scan(&id, &res)
		if err != nil {
			return item, err
		}

		item.Images[id] = res
	}

	rows, err = conn.QueryContext(ctx, "SELECT t.tag FROM tags t JOIN packages p ON p.id = t.pid WHERE p.id = ?", id)
	if err != nil {
		return item, err
	}
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return item, err
		}

		item.Tags = append(item.Tags, tag)
	}

	return item, nil
}

func GetOriginalItemID(ctx context.Context, id int) (int, error) {
	var pid int
	err := conn.QueryRowContext(ctx, "SELECT pid FROM duplicates WHERE id = ?", id).Scan(&pid)
	if err != nil {
		return 0, ErrInvalidID
	}

	return pid, nil
}
