package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hsmtkk/studious-octo-memory/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Handler interface {
	Index(c echo.Context) error
	LoginGet(c echo.Context) error
	LoginPost(c echo.Context) error
	Logout(c echo.Context) error
	HomeGet(c echo.Context) error
	HomePost(c echo.Context) error
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
	user, err := h.requireLogin(c)
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

const sessionName = "ytboard-session"

func (h *handlerImpl) requireLogin(c echo.Context) (model.User, error) {
	ses, _ := session.Get(sessionName, c)
	if ses.Values["login"] == nil || !ses.Values["login"].(bool) {
		c.Redirect(http.StatusFound, "/login")
	}
	ac := ""
	if ses.Values["account"] != nil {
		ac = ses.Values["account"].(string)
	}
	var user model.User
	if err := h.db.Where("account = ?", ac).First(&user).Error; err != nil {
		return user, fmt.Errorf("user query failed; %w", err)
	}
	return user, nil
}

type loginParam struct {
	Title   string
	Message string
	Account string
}

func (h *handlerImpl) LoginGet(c echo.Context) error {
	item := loginParam{
		Title:   "Login",
		Message: "type your account & password:",
		Account: "",
	}
	return c.Render(http.StatusOK, "login", item)
}

func (h *handlerImpl) LoginPost(c echo.Context) error {
	item := loginParam{
		Title:   "Login",
		Message: "type your account & password:",
		Account: "",
	}
	usr := c.FormValue("account")
	pass := c.FormValue("pass")
	item.Account = usr
	var re int64
	var user model.User
	h.db.Where("account = ? AND password = ?", usr, pass).Find(&user).Count(&re)
	if re <= 0 {
		item.Message = "Wrong account or password"
		item.Account = ""
		return c.Render(http.StatusOK, "login", item)
	}

	ses, _ := session.Get(sessionName, c)
	ses.Values["login"] = true
	ses.Values["account"] = usr
	ses.Values["name"] = user.Name
	ses.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/")
}

func (h *handlerImpl) Logout(c echo.Context) error {
	ses, _ := session.Get(sessionName, c)
	ses.Values["login"] = nil
	ses.Values["account"] = nil
	ses.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, "/login")
}

type homeParam struct {
	Title   string
	Message string
	Name    string
	Account string
	Plist   []model.Post
	Glist   []model.Group
}

func (h *handlerImpl) HomeGet(c echo.Context) error {
	user, err := h.requireLogin(c)
	if err != nil {
		return fmt.Errorf("login failed; %w", err)
	}
	var pts []model.Post
	var gps []model.Group
	h.db.Where("user_id = ?", user.ID).Order("created_at desc").Limit(10).Find(&pts)
	h.db.Where("user_id = ?", user.ID).Order("created_at desc").Limit(10).Find(&gps)
	itm := homeParam{
		Title:   "Home",
		Message: fmt.Sprintf("User account=%s", user.Account),
		Name:    user.Name,
		Account: user.Account,
		Plist:   pts,
		Glist:   gps,
	}
	return c.Render(http.StatusOK, "home", itm)
}

func (h *handlerImpl) HomePost(c echo.Context) error {
	user, err := h.requireLogin(c)
	if err != nil {
		return fmt.Errorf("login failed; %w", err)
	}

	switch c.FormValue("form") {
	case "post":
		ad := c.FormValue("address")
		ad = strings.TrimSpace(ad)
		if strings.HasPrefix(ad, "https://youtu.be/") {
			ad = strings.TrimPrefix(ad, "https://youtu.be/")
		}
		pt := model.Post{
			UserID:  int(user.Model.ID),
			Address: ad,
			Message: c.FormValue("message"),
		}
		h.db.Create(&pt)
	case "group":
		gp := model.Group{
			UserID:  int(user.Model.ID),
			Name:    c.FormValue("name"),
			Message: c.FormValue("message"),
		}
		h.db.Create(&gp)
	}

	var pts []model.Post
	var gps []model.Group
	h.db.Where("user_id = ?", user.ID).Order("created_at desc").Limit(10).Find(&pts)
	h.db.Where("user_id = ?", user.ID).Order("created_at desc").Limit(10).Find(&gps)
	itm := homeParam{
		Title:   "Home",
		Message: fmt.Sprintf("User account=%s", user.Account),
		Name:    user.Name,
		Account: user.Account,
		Plist:   pts,
		Glist:   gps,
	}
	return c.Render(http.StatusOK, "home", itm)
}
