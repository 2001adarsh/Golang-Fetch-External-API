package main

import (
	"io"
	"log"
	"nhlpool/nhlApi"
	"os"
	"sync"
	"time"
)

func main() {
	//help benchmarking the request time.
	now := time.Now()

	rosterFile, err := os.OpenFile("rosters.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error in opening the rosters.txt: %v", err)
	}
	defer rosterFile.Close()
	wrt := io.MultiWriter(os.Stdout, rosterFile)
	log.SetOutput(wrt)

	teams, err := nhlApi.GetAllTeams()
	if err != nil {
		log.Fatalf("error while getting teams. %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(teams))

	//unbuffered channels
	results := make(chan []nhlApi.Roster)

	for _, team := range teams {

		go func(team nhlApi.Team) {
			roster, err := nhlApi.GetRosters(team.ID)
			if err != nil {
				log.Fatalf("error in getting roster %v", err)
			}
			results <- roster
			wg.Done() //decrement by 1. one team is done.
		}(team)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	display(results)
	log.Printf("took %v", time.Now().Sub(now).String())

}

func display(results chan []nhlApi.Roster) {
	for r := range results {
		for _, ros := range r {
			log.Println("---------------------")
			log.Printf("Id: %d\n", ros.Person.ID)
			log.Printf("Name: %s\n", ros.Person.FullName)
			log.Printf("Postion: %s\n", ros.Position.Abbreviation)
			log.Printf("Jersey Number: %s\n", ros.JerseyNumber)
			log.Println("---------------------")
		}
	}
}
