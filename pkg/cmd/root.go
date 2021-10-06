package cmd

import (
	"github.com/bal3000/BalStreamerV3/pkg/chromecast"
	"github.com/bal3000/BalStreamerV3/pkg/livestream"
	"github.com/spf13/cobra"
)

func NewCmdRoot(l livestream.LiveStreamer, c chromecast.Chromecaster) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bs <command> [flags]",
		Short:   "BalStreamer CLI",
		Long:    "Control the bal streamer from command line",
		Example: "bs autoplay liverpool -n england",
	}

	cmd.AddCommand(NewAutoPlayCmd(l, c))

	return cmd
}
