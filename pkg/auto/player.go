package auto

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/chromecast"
	"github.com/bal3000/BalStreamerV3/pkg/livestream"
)

type AutoPlayer struct {
	LiveStreamer    livestream.LiveStreamer
	Chromecaster    chromecast.Chromecaster
	SportType       string
	Team            string
	BroadcastNation string
	Chromecast      string
}

func (player *AutoPlayer) ScheduleFixture() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()

	lf, err := player.LiveStreamer.GetLiveFixtures(ctx, player.SportType, now.Format("2006-01-02"), now.Add(24*time.Hour).Format("2006-01-02"), false)
	if err != nil {
		return err
	}

	var fixture []livestream.LiveFixtures

	// filter out only fixtures with my team
	for _, f := range lf {
		if strings.Contains(strings.ToUpper(f.Title), strings.ToUpper(player.Team)) {
			fixture = append(fixture, f)
		}
	}

	if len(fixture) == 0 {
		return fmt.Errorf("no fixtures found for %s", player.Team)
	}

	// see if a fixture with my nation exists
	selectedF := fixture[0]
	for _, f := range fixture {
		if f.BroadcastNationName == player.BroadcastNation {
			selectedF = f
			break
		}
	}

	go func(fixture livestream.LiveFixtures) {
		s := make(chan string, 1)
		ctx, cancel := context.WithCancel(context.Background())

		go player.checkFixture(ctx, 10*time.Minute, fixture, s)

		stream := <-s
		cancel()

		// send event with stream
		player.Chromecaster.CastStream(context.Background(), "chromecast-key", chromecast.StreamToCast{
			Fixture:    fixture.Title,
			StreamURL:  stream,
			Chromecast: player.Chromecast,
		})
	}(selectedF)

	return nil
}

func (player *AutoPlayer) checkFixture(ctx context.Context, interval time.Duration, fixture livestream.LiveFixtures, stream chan<- string) {
	timer := time.NewTimer(interval)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			start, err := time.Parse("2006-01-02", fixture.UtcStart)
			if err != nil {
				log.Printf("error occured converting start date: %w", err)
				return
			}

			if time.Now().Before(start) {
				_ = timer.Reset(interval)
			}

			s := getFixture(player.LiveStreamer, fixture.TimerID)
			if s.RTMP == "" {
				_ = timer.Reset(interval)
			}

			stream <- s.RTMP
			close(stream)
		}
	}
}

func getFixture(liveStreamer livestream.LiveStreamer, timerId string) livestream.Streams {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	streams, err := liveStreamer.GetStreams(ctx, timerId)
	if err != nil {
		log.Printf("error occured getting streams %v\n", err)
		return livestream.Streams{}
	}

	return streams
}
