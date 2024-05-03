package utils

import (
	"streaming/initializers"
	"time"

	"github.com/livekit/protocol/auth"
)

func CreateToken(isStreamer bool, roomId, userId, userName, userAvatar string, config *initializers.Config) (string, error) {

	LiveKitAPISecret := config.LiveKit.APISecret
	LiveKitAPIKey := config.LiveKit.APIKey
	trueVal := true

	at := auth.NewAccessToken(LiveKitAPIKey, LiveKitAPISecret)
	grant := &auth.VideoGrant{
		RoomJoin:       true,
		Room:           roomId,
		CanPublish:     &isStreamer,
		CanPublishData: &trueVal,
	}
	at.AddGrant(grant).
		SetIdentity(userId).
		SetValidFor(6 * time.Hour).
		SetMetadata(userAvatar).
		SetName(userName)

	return at.ToJWT()

}
