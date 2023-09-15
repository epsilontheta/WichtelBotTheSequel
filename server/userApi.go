package server

import (
	"lommix/wichtelbot/server/components"
	"lommix/wichtelbot/server/store"
	"net/http"
)

// ----------------------------------
// User endpoint
func (app *AppState) User(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		err := UserGet(app, writer, request)
		if err != nil {
			println(err.Error())
			http.Error(writer, "forbidden", http.StatusForbidden)
		}
	case http.MethodPut:
		err := userPut(app, writer, request)
		if err != nil {
			println(err.Error())
			http.Error(writer, "forbidden", http.StatusForbidden)
		}
	default:
		http.Error(writer, "forbidden", http.StatusMethodNotAllowed)
		return
	}
}

type updateForm struct {
	Notice    string
	ExcludeId int
}

func userPut(app *AppState, writer http.ResponseWriter, request *http.Request) error {
	form := &updateForm{}
	err := components.FromFormData(request, form)
	if err != nil {
		println(err.Error())
		return err
	}

	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		return err
	}

	user.Notice = form.Notice
	user.ExcludeId = int64(form.ExcludeId)

	err = user.Update(app.Db)
	if err != nil {
		return err
	}

	err = app.Templates.Render(writer, "user", app.defaultContext(request))
	if err != nil {
		return err
	}

	return nil
}

func UserGet(app *AppState, writer http.ResponseWriter, request *http.Request) error {
	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		return err
	}

	if user.PartnerId != 0 {
		partner, err := store.FindUserById(app.Db, user.PartnerId)
		if err == nil {
			user.Partner = &partner
		}
	}

	err = app.Templates.Render(writer, "user", user)

	if err != nil {
		return err
	}
	return nil
}
