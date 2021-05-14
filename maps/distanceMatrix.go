package maps

import (
	"context"
	"fmt"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
)


func DistanceMatrix(originId string, destinationId string) {

	var client *maps.Client
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("new client %v\n", err)
	}
	if client == nil {
		log.Fatalln("no client")
	}

	r := &maps.DistanceMatrixRequest{
		ArrivalTime:   "1615194000",
	}

	r.Origins = []string{fmt.Sprintf("place_id:%s", originId)}
	r.Destinations = []string{fmt.Sprintf("place_id:%s", destinationId)}

	r.Mode = maps.TravelModeTransit
	r.TransitMode = []maps.TransitMode{
		maps.TransitModeBus,
		maps.TransitModeSubway,
	}

	resp, err := client.DistanceMatrix(context.Background(), r)

	_, _ = pretty.Println(resp)
}

