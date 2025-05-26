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
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MetadataItem struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Images            []string `json:"images"`
	ReuploaderName    string   `json:"reuploader_name"`
	ReuploaderProfile string   `json:"reuploader_profile"`
	Reuploaded        int      `json:"reuploaded"` // unix time
	Updated           int      `json:"updated"`    // unix time
	Downloads         int      `json:"downloads"`
	Views             int      `json:"views"`
	Size              any      `json:"size"` // use GetSize
	Tags              []string `json:"tags"`
	Description       string   `json:"description"`
	OrigUploader      string   `json:"orig_uploader"`
	OrigUploadDate    int      `json:"orig_uploaddate"`
	PageURL           string   `json:"pageurl"`
	DLIndirect        string   `json:"dl_indirect"`
	DLDirect          string   `json:"dl_direct"`
}

var ErrUnknownID = errors.New("unknown id")

// returns metadata from id
func GetMetadataFromID(id int) (MetadataItem, error) {
	for _, m := range metadata {
		mid, err := m.GetID()
		if err != nil {
			continue
		}

		if mid != id {
			continue
		}

		return m, nil
	}

	return MetadataItem{}, ErrUnknownID
}

func (m MetadataItem) GetID() (int, error) {
	i, err := strconv.Atoi(m.ID)
	if err != nil {
		return 0, err
	}

	return i, nil
}

var ErrUnhandledType = errors.New("unhandled type")

func (m MetadataItem) GetSize() (float64, error) {
	switch m.Size.(type) {
	case float64:
		return m.Size.(float64), nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSuffix(m.Size.(string), " B"), 64)
		if err != nil {
			return 0, err
		}

		return f, nil
	}

	return 0, ErrUnhandledType
}

func (m MetadataItem) GetPrettySize() (string, error) {
	size, err := m.GetSize()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.02f MB", size/1024/1024), nil
}

func (m MetadataItem) GetUploadTime() time.Time {
	if m.OrigUploadDate != 00 {
		return time.Unix(int64(m.OrigUploadDate), 0)
	}

	return time.Unix(int64(m.Reuploaded), 0)
}

func (m MetadataItem) GetUploaderName() string {
	if m.OrigUploader != "" {
		return m.OrigUploader
	}
	if m.ReuploaderName != "garrysmod.org" {
		return m.ReuploaderName
	}

	return "Unknown"
}

func (m MetadataItem) GetUploaderID() (int, error) {
	id, err := strconv.Atoi(strings.TrimPrefix(m.ReuploaderProfile, "https://garrysmods.org/profile/"))
	if err != nil {
		return 0, err
	}

	return id, nil
}
