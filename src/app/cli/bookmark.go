package cli

import (
	"klog/app"
	"klog/app/cli/lib"
)

type Bookmark struct {
	//List   BookmarkList   `cmd name:"list" group:"Bookmark" help:"Show all available bookmarks"` // TODO
	Set    BookmarkSet    `cmd name:"set" group:"Bookmark" help:"Set bookmark to a file"`
	Unset  BookmarkUnset  `cmd name:"unset" group:"Bookmark" help:"Clear current bookmark"`
	//Rename BookmarkRename `cmd name:"rename" group"Bookmark" help:"Change the name of a bookmark"` // TODO
	Edit   BookmarkEdit   `cmd name:"edit" group:"Bookmark" help:"Open bookmark in your $EDITOR"`
}

func (opt *Bookmark) Help() string {
	return `With bookmarks you can make klog always read from a default file, in case you don’t specify one explicitly.

This is handy in case you always use the same file.
You can then interact with it regardless of your current working directory.`
}

type BookmarkSet struct {
	File     string `arg type:"existingfile" help:".klg source file"`
	Bookmark string `arg help:"The name for the bookmark"`
	lib.QuietArgs
}

func (args *BookmarkSet) Run(ctx app.Context) error {
	//err := ctx.BookmarksWrite(args.File)
	//if err != nil {
	//	return err
	//}
	if !args.Quiet {
		ctx.Print("Bookmarked file ")
	}
	ctx.Print(args.File + "\n")
	return nil
}

type BookmarkEdit struct{}

func (args *BookmarkEdit) Run(ctx app.Context) error {
	b, appErr := ctx.Bookmark()
	if appErr != nil {
		return appErr
	}
	return ctx.OpenInEditor(b.Path)
}

type BookmarkUnset struct{}

func (args *BookmarkUnset) Run(ctx app.Context) error {
	err := ctx.UnsetBookmark()
	if err != nil {
		return err
	}
	ctx.Print("Cleared bookmark\n")
	return nil
}
