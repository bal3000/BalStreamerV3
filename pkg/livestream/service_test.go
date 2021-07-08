package livestream

import (
	"testing"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/config"
)

var c = config.Configuration{
	LiveStreamURL: "http://test.com",
	APIKey:        "1234",
}

func TestFilterLiveFixtures(t *testing.T) {
	fixtures := []LiveFixtures{
		{
			StateName:            "running",
			UtcStart:             time.Now().Add(-10 * time.Minute).Format("2006-01-02T15:04:05"),
			UtcEnd:               time.Now().Add(2 * time.Hour).Format("2006-01-02T15:04:05"),
			Title:                "Cerezo Osaka vs Guangzhou FC",
			EventID:              "1e9i6jwhskjzdhk1h9o12776c",
			ContentTypeName:      "Soccer",
			TimerID:              "129545",
			IsPrimary:            "true",
			BroadcastChannelName: "C More Sport 1",
			BroadcastNationName:  "Finland",
			SourceTypeName:       "Sat-Receiver",
		},
		{
			StateName:            "upcoming",
			UtcStart:             time.Now().Add(2 * time.Hour).Format("2006-01-02T15:04:05"),
			UtcEnd:               time.Now().Add(4 * time.Hour).Format("2006-01-02T15:04:05"),
			Title:                "BU2: Italy vs Spain",
			EventID:              "8w1xokiua6sq5x6j0xy16t6y2",
			ContentTypeName:      "Soccer",
			TimerID:              "129534",
			IsPrimary:            "false",
			BroadcastChannelName: "ESPN",
			BroadcastNationName:  "USA",
			SourceTypeName:       "Sat-Receiver",
		},
		{
			StateName:            "upcoming",
			UtcStart:             time.Now().Add(2 * time.Hour).Format("2006-01-02T15:04:05"),
			UtcEnd:               time.Now().Add(4 * time.Hour).Format("2006-01-02T15:04:05"),
			Title:                "Japan vs South Korea",
			EventID:              "8w1xdsadasa6sq5x6j0xy16t6y2",
			ContentTypeName:      "Soccer",
			TimerID:              "129532",
			IsPrimary:            "false",
			BroadcastChannelName: "Sky Sports",
			BroadcastNationName:  "England",
			SourceTypeName:       "Sat-Receiver",
		},
	}

	expected := []LiveFixtures{
		{
			StateName:            "running",
			UtcStart:             time.Now().Add(-10 * time.Minute).Format("2006-01-02T15:04:05"),
			UtcEnd:               time.Now().Add(2 * time.Hour).Format("2006-01-02T15:04:05"),
			Title:                "Cerezo Osaka vs Guangzhou FC",
			EventID:              "1e9i6jwhskjzdhk1h9o12776c",
			ContentTypeName:      "Soccer",
			TimerID:              "129545",
			IsPrimary:            "true",
			BroadcastChannelName: "C More Sport 1",
			BroadcastNationName:  "Finland",
			SourceTypeName:       "Sat-Receiver",
		},
	}

	s := NewService(c)
	lf, err := s.FilterLiveFixtures(fixtures)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if len(lf) != len(expected) {
		t.Errorf("expected result length of %v got %v", len(expected), len(lf))
	}

	for i, result := range lf {
		if result.Title != expected[i].Title {
			t.Errorf("expected title of %s got %s", expected[i].Title, result.Title)
		}
	}
}
