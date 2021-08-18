package app

import "os"

type legacyBookmark struct {
	basePath    string
	symlinkName string
}

func (l *legacyBookmark) originPath() string {
	return l.basePath + "bookmark.klg"
}

func (l *legacyBookmark) get() *BookmarksCollection {
	path, err := os.Readlink(l.originPath())
	if err != nil {
		return nil
	}
	bookmarksCollection := NewBookmarksCollection()
	bookmarksCollection.Set(DEFAULT, path)
	return bookmarksCollection
}

func (l *legacyBookmark) cleanup() {
	_ = os.Remove(l.originPath()) // Disregard error, since it’s a legacy file anyway
}
