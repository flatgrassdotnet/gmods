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

import "context"

func GetPopularTags(ctx context.Context) ([]string, error) {
	var tags []string
	rows, err := conn.QueryContext(ctx, "SELECT tag FROM tags GROUP BY tag ORDER BY COUNT(*) DESC, tag ASC LIMIT 100")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
