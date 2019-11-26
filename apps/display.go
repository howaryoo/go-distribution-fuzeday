package apps

import (
	"fmt"
	"go-distribution-fuzeday/messaging"
	"go-distribution-fuzeday/models"
	"net/http"
	"sync"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "The Football Game Of Tikal!")
}

func displayGamefield(gamefield *models.GameField) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bb, _ := gamefield.MarshalJSON()
		s := string(bb)
		fmt.Fprint(w, s)
	}
}

func displayGamefieldStr(gamefield models.GameField) string {
	bb, _ := gamefield.MarshalJSON()
	return string(bb)
}

func LaunchDisplay(port int, externalWaitGroup *sync.WaitGroup) {
	displayInput := getDisplayInputChannel()

	gameField := models.NewGameField()

	go receiveDisplayChannelUpdates(displayInput, gameField)

	// HTTP Server
	http.HandleFunc("/", myHandler)
	//	2. requests to "/display" should return a json representation of the updated gameField
	http.HandleFunc("/display", displayGamefield(gameField))
	//	3. requests to "/client/" should return static files from directory "display_client". Use http.FileServer...
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("display_client"))))
	// 	------
	//	1. launch HTTP server here on 8080
	// 	Tip: use http.HandleFunc and http.ListenAndServe
	panic(http.ListenAndServe(":8080", nil))

	externalWaitGroup.Wait()
	if externalWaitGroup != nil {
		fmt.Println("externalWaitGroup.Done!")
		externalWaitGroup.Done()
	}
}

	// Game Field updater
	//	------
	//	Tip: use iteration over channel range
func receiveDisplayChannelUpdates(displayInput chan *models.DisplayStatus, gamefield *models.GameField) {
	//	1. iterate over display channel
	for {
		select {
		case msg := <-displayInput:
			//	2. update gamefield on each consumed value
			gamefield.Update(msg)
		default:
			//fmt.Println("no message received")
		}
	}
}

func getDisplayInputChannel() chan *models.DisplayStatus {
	//TODO Challenge (2):
	//  get []byte input channel from messaging,
	//  create an internal goroutine that consumes messages from it,
	//  de-serialize them to return type and populates return DIRECTIONAL channel
	return messaging.GlobalDisplayChannel
}
