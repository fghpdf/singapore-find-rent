package maps

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

const apiKey = "AIzaSyAVIIF1Y4Q0FN2S2zFI23rIszkCvDqdycM"

var (
	input        = flag.String("input", "", "The text input specifying which place to search for (for example, a name, address, or phone number).")
)

func FindPlaceFromText() string {
	flag.Parse()

	var client *maps.Client
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("parse %v\n", err)
	}

	req := &maps.FindPlaceFromTextRequest{
		Input:     *input,
		InputType: maps.FindPlaceFromTextInputTypeTextQuery,
	}

	fields, err := parseFields("formatted_address,name,place_id")
	if err != nil {
		log.Fatalf("parse %v\n", err)
	}
	req.Fields = fields

	resp, err := client.FindPlaceFromText(context.Background(), req)

	_, _ = pretty.Println(resp)

	if len(resp.Candidates) != 0 {
		return resp.Candidates[0].PlaceID
	}

	return ""
}

func parseFields(fields string) ([]maps.PlaceSearchFieldMask, error) {
	var res []maps.PlaceSearchFieldMask
	for _, s := range strings.Split(fields, ",") {
		f, err := maps.ParsePlaceSearchFieldMask(s)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}
