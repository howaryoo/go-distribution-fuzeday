package apps

import (
	"fmt"
	"go-distribution-fuzeday/messaging"
	"go-distribution-fuzeday/models"
	"net/http"
	"sync"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "The Football Game Of Tikal!")
}

func displayGamefield(gamefield *models.GameField) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bb, _ := gamefield.MarshalJSON()
		s := string(bb)
		//fmt.Println("gamefieled is ", s)
		fmt.Fprint(w, s)
	}
}

func displayGamefieldStr(gamefield models.GameField) string {
	bb, _ := gamefield.MarshalJSON()
	return string(bb)
}

func LaunchDisplay(port int, externalWaitGroup *sync.WaitGroup) {
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++ Launch Display Started")
	displayInput := getDisplayInputChannel()

	gameField := models.NewGameField()

	go my_iterate(displayInput, gameField)

	//testDisplayInput(displayInput)	// sends 1 DisplayStatus into the channel

	// HTTP Server
	//TODO Challenge (4):
	//	1. launch HTTP server here on 8080

	http.HandleFunc("/", myHandler)

	//	2. requests to "/display" should return a json representation of the updated gameField
	http.HandleFunc("/display", displayGamefield(gameField))
	//	3. requests to "/client/" should return static files from directory "display_client". Use http.FileServer...
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("display_client"))))
	// 	------
	// 	Tip: use http.HandleFunc and http.ListenAndServe
	panic(http.ListenAndServe(":8080", nil))

	// Game Field updater
	//TODO Challenge (4):
	//	1. iterate over display channel
	//	2. update gamefield on each consumed value
	//	------
	//	Tip: use iteration over channel range
	//go my_iterate(displayInput)

	//displayInput = displayInput // only to prevent "unused variable error", remove after implementation
	//gameField = gameField       // only to prevent "unused variable error", remove after implementation
	externalWaitGroup.Wait()
	if externalWaitGroup != nil {
		fmt.Println("externalWaitGroup.Done!")
		externalWaitGroup.Done()
	}
}

func testDisplayInput(displayInput chan *models.DisplayStatus) {
	displayInput <- &models.DisplayStatus{
		ItemID:      "1",
		ItemLabel:   "2",
		ItemType:    "3",
		TeamID:      "4",
		X:           0,
		Y:           0,
		Z:           0,
		LastUpdated: time.Time{},
	}
}

func my_iterate(displayInput chan *models.DisplayStatus, gamefield *models.GameField) {
	fmt.Println("++++++++++++ my_iterate started")
	for {
		select {
		case msg := <-displayInput:
			//fmt.Println("++++++++++++ received message ", msg)
			gamefield.Update(msg)
			//fmt.Println(displayGamefieldStr(gamefield))
		default:
			//fmt.Println("no message received")
		}
	}

	//for item := range displayInput {
	//	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	//	fmt.Println("+++")
	//	fmt.Println(item)
	//	fmt.Println("+++")
	//	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	//}
}

func getDisplayInputChannel() chan *models.DisplayStatus {
	//TODO Challenge (2):
	//  get []byte input channel from messaging,
	//  create an internal goroutine that consumes messages from it,
	//  de-serialize them to return type and populates return DIRECTIONAL channel
	return messaging.GlobalDisplayChannel
}
