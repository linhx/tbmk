package bookmark

import (
	"fmt"
	"os"

	common "linhx.com/tbmk/common"

	"github.com/sahilm/fuzzy"
	"github.com/sonyarouje/simdb"
	"github.com/google/uuid"
)

type BookmarkItem struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Command string `json:"command"`
}

func (bmki BookmarkItem) ID() (jsonField string, value interface{}) {
	value = bmki.Id
	jsonField = "id"
	return
}

type BookmarkRepo struct {
	db *simdb.Driver
}

func NewBookmarkRepo() (*BookmarkRepo, error) {
	repo := new(BookmarkRepo)
	driver, err := simdb.New(os.Getenv("TBMK_DATA_DIR"))
	repo.db = driver
	return repo, err
}

func (repo *BookmarkRepo) createNewBookmarkItemId() (string) {
	return uuid.New().String()
}

func (repo *BookmarkRepo) getAllBookmarkItems() ([]BookmarkItem, error) {
	var bmkis []BookmarkItem
	err := repo.db.Open(BookmarkItem{}).Get().AsEntity(&bmkis)
	return bmkis, err
}

func (repo *BookmarkRepo) save(title string, command string, override bool) (BookmarkItem, error) {
	var bmki BookmarkItem
	var err error
	var newId string
	_ = repo.db.Open(BookmarkItem{}).Where("title", "=", title).First().AsEntity(&bmki)
	if len(bmki.Id) > 0 {
		if override {
			bmki.Command = command
			err = repo.db.Update(bmki)
			return bmki, err
		} else {
			return bmki, common.NewDuplicateBmkiError(fmt.Sprintf("Already exist title '%s'", title), bmki.Id)
		}
	}
	_ = repo.db.Open(BookmarkItem{}).Where("command", "=", command).First().AsEntity(&bmki)
	if len(bmki.Id) > 0 {
		if override {
			bmki.Title = title
			err = repo.db.Update(bmki)
			return bmki, err
		} else {
			return bmki, common.NewDuplicateBmkiError(fmt.Sprintf("Already exist command '%s'", command), bmki.Id)
		}
	}

	newId = repo.createNewBookmarkItemId()
	bmki = BookmarkItem{
		Id:      newId,
		Title:   title,
		Command: command,
	}
	err = repo.db.Open(BookmarkItem{}).Upsert(bmki)
	return bmki, err
}

func (repo *BookmarkRepo) remove(id string) error {
	bmkiToDelete := BookmarkItem{
		Id: id,
	}
	err := repo.db.Delete(bmkiToDelete)
	return err
}

type Bookmark struct {
	repo       *BookmarkRepo
	cacheBmkis *[]BookmarkItem
}

func NewBookmark() (*Bookmark, error) {
	bmk := new(Bookmark)
	repo, err := NewBookmarkRepo()
	bmk.repo = repo
	return bmk, err
}

func (bmk *Bookmark) Save(title string, command string, override bool) (BookmarkItem, error) {
	return bmk.repo.save(title, command, override)
}

func (bmk *Bookmark) Remove(id string) error {
	err := bmk.repo.remove(id)

	var cacheBmkis = *(bmk.cacheBmkis)
	var index int = -1
	for i, n := range cacheBmkis {
		if n.Id == id {
			index = i
			break
		}
	}
	if index > -1 {
		cacheBmkis = append(cacheBmkis[:index], cacheBmkis[index+1:]...)
		bmk.cacheBmkis = &cacheBmkis
	}
	return err
}

type MatchedItem struct {
	Id           string
	Title        string
	Command      string
	MatchTitle   fuzzy.Match
	MatchCommand fuzzy.Match
}

func FindIndex(maches []MatchedItem, id string) int {
	for i, n := range maches {
		if n.Id == id {
			return i
		}
	}
	return -1
}

func (bmk *Bookmark) Search(query string) ([]MatchedItem, error) {
	var matchesBmkis []MatchedItem
	if bmk.cacheBmkis == nil {
		allBmkis, err := bmk.repo.getAllBookmarkItems()
		bmk.cacheBmkis = &allBmkis
		if err != nil {
			return matchesBmkis, err
		}
	}
	var allBmkis = *bmk.cacheBmkis
	if len(query) == 0 {
		for _, bmki := range allBmkis {
			matchesBmkis = append(matchesBmkis, MatchedItem{
				Id:      bmki.Id,
				Title:   bmki.Title,
				Command: bmki.Command,
			})
		}
		return matchesBmkis, nil
	}
	var allBmkisTitle []string
	var allBmkisCommand []string
	for _, bmki := range allBmkis {
		allBmkisTitle = append(allBmkisTitle, bmki.Title)
		allBmkisCommand = append(allBmkisCommand, bmki.Command)
	}

	matchesTitle := fuzzy.Find(query, allBmkisTitle)
	matchesCommand := fuzzy.Find(query, allBmkisCommand)

	for _, match := range matchesTitle {
		_bmki := allBmkis[match.Index]
		matchesBmkis = append(matchesBmkis, MatchedItem{
			Id:         _bmki.Id,
			Title:      _bmki.Title,
			Command:    _bmki.Command,
			MatchTitle: match,
		})
	}

	for _, match := range matchesCommand {
		_bmki := allBmkis[match.Index]
		existBmkiI := FindIndex(matchesBmkis, _bmki.Id)
		if existBmkiI > -1 {
			existBmki := matchesBmkis[existBmkiI]
			existBmki.MatchCommand = match
			matchesBmkis[existBmkiI] = existBmki
			continue
		} else {
			matchesBmkis = append(matchesBmkis, MatchedItem{
				Id:           _bmki.Id,
				Title:        _bmki.Title,
				Command:      _bmki.Command,
				MatchCommand: match,
			})
		}
	}

	return matchesBmkis, nil
}
