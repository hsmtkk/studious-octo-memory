package handle

import (
	"fmt"
	"net/http"

	"github.com/hsmtkk/studious-octo-memory/login"
	"github.com/hsmtkk/studious-octo-memory/model"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Handler interface {
	Index(c echo.Context) error
}

type handlerImpl struct {
	db *gorm.DB
}

func New(dbPath string) (Handler, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database; %s; %w", dbPath, err)
	}
	return &handlerImpl{db: db}, nil
}

type indexParam struct {
	Title   string
	Message string
	Name    string
	Account string
	Plist   []model.Post
	Glist   []model.Group
}

func (h *handlerImpl) Index(c echo.Context) error {
	user, err := login.RequireLogin(c)
	if err != nil {
		return fmt.Errorf("login failed; %w", err)
	}

	var posts []model.Post
	if err := h.db.Where("group_id > 0").Order("created_at desc").Limit(10).Find(&posts).Error; err != nil {
		return fmt.Errorf("post query failed; %w", err)
	}
	var groups []model.Group
	if err := h.db.Order("created_at desc").Limit(10).Find(&groups).Error; err != nil {
		return fmt.Errorf("group query failed; %w", err)
	}

	item := indexParam{
		Title:   "Index",
		Message: "This is top page",
		Name:    user.Name,
		Account: user.Account,
		Plist:   posts,
		Glist:   groups,
	}

	return c.Render(http.StatusOK, "index", item)
}
