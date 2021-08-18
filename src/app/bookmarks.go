package app

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
)

var DEFAULT = NewAlias("default")

type BookmarkAlias string

type Bookmark struct {
	Alias  BookmarkAlias
	Target *File
}

type BookmarksCollection map[BookmarkAlias]*File

type BookmarkJson struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}

func ParseBookmarks(jsonText string) (*BookmarksCollection, Error) {
	var data []BookmarkJson
	err := json.Unmarshal([]byte(jsonText), &data)
	if err != nil {
		return nil, NewErrorWithCode(
			BOOKMARK_ACCESS_ERROR,
			"Cannot read bookmark file",
			"The file ~/.klog/bookmarks.json is not correctly formatted",
			err,
		)
	}
	collection := NewBookmarksCollection()
	for _, d := range data {
		collection.Set(NewAlias(d.Alias), d.Path)
	}
	return collection, nil
}

func NewBookmarksCollection() *BookmarksCollection {
	var collection BookmarksCollection
	collection = make(map[BookmarkAlias]*File)
	return &collection
}

func (c *BookmarksCollection) ToJson() string {
	data := make([]BookmarkJson, len(*c))
	i := 0
	for alias, file := range *c {
		data[i] = BookmarkJson{
			Path:  file.Path,
			Alias: alias.RawName(),
		}
		i++
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Alias < data[j].Alias
	})
	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	err := enc.Encode(&data)
	if err != nil {
		// This can/should never happen
		panic(err)
	}
	return buffer.String()
}

func (c *BookmarksCollection) Count() int {
	return len(*c)
}

func (c *BookmarksCollection) All() []Bookmark {
	var result []Bookmark
	for name, file := range *c {
		result = append(result, Bookmark{name, file})
	}
	return result
}

func (c *BookmarksCollection) Lookup(alias BookmarkAlias) (*Bookmark, *AppError) {
	file := (*c)[alias]
	if file == nil {
		return nil, &AppError{
			code:     BOOKMARK_NOT_SET,
			message:  "No such bookmark",
			details:  "There is no bookmark with that alias",
			original: nil,
		}
	}
	return &Bookmark{alias, file}, nil
}

func (c *BookmarksCollection) GetDefault() *Bookmark {
	defaultBookmarkPath := (*c)[DEFAULT]
	if defaultBookmarkPath == nil {
		return nil
	}
	return &Bookmark{DEFAULT, defaultBookmarkPath}
}

func (c *BookmarksCollection) Set(alias BookmarkAlias, path string) {
	(*c)[alias] = newFile(path)
}

func (c *BookmarksCollection) Rename(oldAlias BookmarkAlias, newAlias BookmarkAlias) *AppError {
	bookmark, err := c.Lookup(oldAlias)
	if err != nil {
		return err
	}
	(*c)[newAlias] = bookmark.Target
	c.Unset(oldAlias)
	return nil
}

func (c *BookmarksCollection) Unset(name BookmarkAlias) {
	delete(*c, name)
}

func NewAlias(name string) BookmarkAlias {
	return BookmarkAlias(strings.TrimPrefix(name, "@"))
}

func (a BookmarkAlias) PrettyName() string {
	return "@" + a.RawName()
}

func (a BookmarkAlias) RawName() string {
	return string(a)
}
