package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsesEmptyJsonIntoEmptyBookmarkCollection(t *testing.T) {
	json := `[]`
	bookmarks, err := ParseBookmarks(json)
	require.Nil(t, err)
	assert.Equal(t, 0, bookmarks.Count())
}

func TestParsingFailsIfJsonIsInvalidOrMalformed(t *testing.T) {
	for _, json := range []string {
		``,
		`[{"foo": `,
		`{"alias": "foo", "xyz": true}`,
	} {
		bookmarks, err := ParseBookmarks(json)
		require.Error(t, err)
		assert.Nil(t, bookmarks)
		assert.Equal(t, err.Code(), BOOKMARK_ACCESS_ERROR)
	}
}

func TestParsesJsonIntoBookmarkCollection(t *testing.T) {
	json := `[
	{"alias": "foo", "path": "~/foo.klg"},
	{"alias": "bar", "path": "~/bar.klg"}
]`
	bookmarks, err := ParseBookmarks(json)
	require.Nil(t, err)
	assert.Equal(t, 2, bookmarks.Count())

	fooBookmark, fooErr := bookmarks.Lookup("foo")
	require.Nil(t, fooErr)
	assert.Equal(t, "foo", fooBookmark.Alias)
	assert.Equal(t, "~/foo.klg", fooBookmark.Target.Path)

	barBookmark, barErr := bookmarks.Lookup("bar")
	require.Nil(t, barErr)
	assert.Equal(t, "bar", barBookmark.Alias)
	assert.Equal(t, "~/bar.klg", barBookmark.Target.Path)
}

func TestSerialisesBookmarksCorrectly(t *testing.T) {
	json := `[
	{"alias": "foo", "path": "~/foo.klg"},
	{"alias": "bar", "path": "~/bar.klg"}
]`
	bookmarks, _ := ParseBookmarks(json)
	// Sorted alphabetically (by name) and formatted nicely:
	assert.Equal(t, `[
  {
    "alias": "bar",
    "path": "~/bar.klg"
  },
  {
    "alias": "foo",
    "path": "~/foo.klg"
  }
]
`, bookmarks.ToJson())
}

func TestEmptyDefault(t *testing.T) {
	emptyBookmarks, _ := ParseBookmarks("[]")
	assert.Nil(t, emptyBookmarks.GetDefault())

	noDefault, _ := ParseBookmarks(`[
		{"alias": "foo", "path": "~/foo.klg"},
		{"alias": "bar", "path": "~/bar.klg"}
	]`)
	assert.Nil(t, noDefault.GetDefault())
}

func TestDefaultBookmark(t *testing.T) {
	bookmarks, _ := ParseBookmarks(`[
		{"alias": "default", "path": "~/foo.klg"},
		{"alias": "bar", "path": "~/bar.klg"}
	]`)
	def := bookmarks.GetDefault()
	require.NotNil(t, def)
	assert.Equal(t, "default", def.Alias)
	assert.Equal(t, "~/foo.klg", def.Target.Path)
}

func TestRenameBookmarkFailsIfNotFound(t *testing.T) {
	bookmarks, _ := ParseBookmarks(`[
		{"alias": "foo", "path": "~/foo.klg"},
		{"alias": "bar", "path": "~/bar.klg"}
	]`)
	renameErr := bookmarks.Rename("doesnotexist", "asdf")
	require.Error(t, renameErr)
}

func TestRenameBookmark(t *testing.T) {
	bookmarks, _ := ParseBookmarks(`[
		{"alias": "foo", "path": "~/foo.klg"},
		{"alias": "bar", "path": "~/bar.klg"}
	]`)
	countBefore := bookmarks.Count()
	renameErr := bookmarks.Rename("foo", "asdf")
	require.Nil(t, renameErr)

	asdfBookmark, _ := bookmarks.Lookup("asdf")
	assert.NotNil(t, asdfBookmark)

	fooBookmark, fooErr := bookmarks.Lookup("foo")
	require.Nil(t, fooBookmark)
	assert.Error(t, fooErr)

	assert.Equal(t, countBefore, bookmarks.Count())
}
