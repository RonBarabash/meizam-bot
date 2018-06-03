package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/meizam-bot/controller"
)

func main() {
	messenger := &messenger.Messenger{
		VerifyToken: os.Getenv("VERIFY_TOKEN"),
		AppSecret:   os.Getenv("APP_SECRET"),
		AccessToken: os.Getenv("PAGE_ACCESS_KEY"),
	}

	//dbConnectionParams := providers.DBConnectionParams{
	//	User:     os.Getenv("DB_USER"),
	//	Password: os.Getenv("DB_PASS"),
	//	Address:  os.Getenv("DB_ADDRESS"),
	//	DBName:   os.Getenv("DB_NAME"),
	//}

	//dataProvider, err := providers.NewBotDataProvider(dbConnectionParams)
	//if err != nil {
	//	fmt.Println("could not create suchef server. error: " + err.Error())
	//	return
	//}

	//accountID := int64(1)

	//suchefServer := server.NewSuchefServer(accountID, messenger, dataProvider, dataProvider, dataProvider)
	//fmt.Println("server started successfully")

	//messenger.MessageReceived = suchefServer.BindMessageReceived()
	//messenger.Postback = suchefServer.BindPostbackReceived()
	ctrl := controller.NewController()
	messenger.MessageReceived = ctrl.BindMessageReceived()
	r := mux.NewRouter()
	r.HandleFunc("/webhook", messenger.Handler)
	http.Handle("/", r)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
