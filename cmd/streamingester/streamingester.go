package streamingester

import (
	"fmt"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
)

// NewCmd registers the cobra command.
func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "streamingester",
		Short: "launches a stream ingest node",
		Run:   Start,
	}
}

var que *pubsub.Queue

func Start(_ *cobra.Command, _ []string) {
	log.Info("Starting livestream ingester")

	bindAddress := viper.GetString("stream_ingester.bind_address")

	format.RegisterAll()

	server := &rtmp.Server{
		Addr:          bindAddress,
		HandlePublish: HandlePublish,
		HandlePlay:    handlePlay,
	}

	fmt.Println("Info: Starting the stream server at", bindAddress)
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, "Fatal: Couldn't run the stream server:", err)
	}
}

func HandlePublish(conn *rtmp.Conn) {
	defer conn.Close()

	if que != nil {
		que.Close()
	}
	que = pubsub.NewQueue()

	streams, err := conn.Streams()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Couldn't stream:", err)
		return
	}
	que.WriteHeader(streams)

	fmt.Println("Info: The server has started streaming.")

	if err := avutil.CopyPackets(que, conn); err == io.EOF {
		fmt.Println("Info: The server has stopped streaming.")
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Couldn't stream:", err)
	}
}

func handlePlay(conn *rtmp.Conn) {
	defer conn.Close()
	if que == nil {
		return
	}

	if err := avutil.CopyFile(conn, que.Latest()); err != nil && err != io.EOF {
		fmt.Printf("%+v\n", err)
		fmt.Println("Info: Couldn't serve the stream to a viewer:", err)
	}
}