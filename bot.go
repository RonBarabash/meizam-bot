package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/RonBarabash/meizam-bot/controller"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/RonBarabash/meizam-bot/meizam"
	"github.com/RonBarabash/meizam-bot/providers"
)

func main() {
	messenger := &messenger.Messenger{
		VerifyToken: os.Getenv("VERIFY_TOKEN"),
		AppSecret:   os.Getenv("APP_SECRET"),
		AccessToken: os.Getenv("PAGE_ACCESS_KEY"),
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=MeizamDB", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"))
	m := meizam.NewMeizam(connString)
	messengerProvider := providers.NewFacebookMessengerProvider(messenger)

	ctrl := controller.NewController(m, messengerProvider)
	messenger.MessageReceived = ctrl.BindMessageReceived()
	messenger.Postback = ctrl.BindPostbackReceived()
	messenger.Authentication = ctrl.BindAuthentication()
	r := mux.NewRouter()
	r.HandleFunc("/webhook", messenger.Handler)
	http.Handle("/", r)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

