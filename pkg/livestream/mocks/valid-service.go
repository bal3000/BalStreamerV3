package mocks

import (
	"context"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/livestream"
)

var (
	fixtures = []livestream.LiveFixtures{
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
	rtmpLink = "rtmp://cdn.vops.gcp.xeatre.cloud:5222/liveedge-lowlatency-origin-wza-03/src-3905?wUzz3Tsnestarttime=1625602047&wUzz3Tsneendtime=1625616447&wUzz3Tsnehash=oXr3YM5y4FVOHejFbkHian7-X4vF_KIGwmIBaiP4Qdg=&DVR&wowzadvrplayliststart=20210706185000"
)

type MockService struct{}

func (s MockService) GetFixtureCount() int {
	return len(fixtures)
}

func (s MockService) GetRMTPLink() string {
	return rtmpLink
}

func (s MockService) GetLiveFixtures(ctx context.Context, sportType, fromDate, toDate string, live bool) ([]livestream.LiveFixtures, error) {
	return fixtures, nil
}

func (s MockService) GetStreams(ctx context.Context, timerID string) (livestream.Streams, error) {
	return livestream.Streams{
		RTMP: rtmpLink,
	}, nil
}

func (s MockService) FilterLiveFixtures(fixtures []livestream.LiveFixtures) ([]livestream.LiveFixtures, error) {
	return []livestream.LiveFixtures{
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
	}, nil
}
