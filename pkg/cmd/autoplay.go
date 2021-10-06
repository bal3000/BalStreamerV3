package cmd

import (
	"log"

	"github.com/bal3000/BalStreamerV3/pkg/auto"
	"github.com/bal3000/BalStreamerV3/pkg/chromecast"
	"github.com/bal3000/BalStreamerV3/pkg/livestream"
	"github.com/spf13/cobra"
)

var (
	Nation     string
	Chromecast string
	SportType  string
)

func NewAutoPlayCmd(l livestream.LiveStreamer, cc chromecast.Chromecaster) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "autoplay [team] [flags]",
		Short:   "autoplay a teams fixture",
		Long:    "Schedule for a fixture to play automatically for the given team",
		Example: "bs autoplay liverpool -n united kingdom",
		Args:    cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			player := &auto.AutoPlayer{
				LiveStreamer:    l,
				Chromecaster:    cc,
				Team:            args[0],
				BroadcastNation: Nation,
				Chromecast:      Chromecast,
				SportType:       SportType,
			}

			err := player.ScheduleFixture()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringVarP(&Nation, "nation", "n", "United Kingdom", "Determines which country to get the feed from")
	cmd.Flags().StringVarP(&Chromecast, "chromecast", "c", "Family room TV", "Determines which chromecast to cast to")
	cmd.Flags().StringVarP(&SportType, "sport", "s", "Soccer", "Determines which sport to search by")

	return cmd
}
