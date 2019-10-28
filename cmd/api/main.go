package main

import (
	"fmt"
	"time"

	"github.com/ninnemana/rpc-demo/pkg/vinyltap"
)

func main() {
	a := vinyltap.Album{
		Id:          1,
		Artist:      "Nirvana",
		Title:       "Nevermind",
		ReleaseDate: time.Date(1991, time.September, 24, 0, 0, 0, 0, time.UTC).Unix(),
		Songs: []string{
			"Smells Like Teen Spirit",
			"In Bloom",
			"Come as You Are",
			"Breed",
			"Lithium",
			"Polly",
			"Territorial Pissings",
			"Drain You",
			"Lounge Act",
			"Stay Away",
			"On a Plain",
			"Something in the Way",
			"Endless, Nameless",
		},
	}
	fmt.Printf("%+v\n", a)
}
