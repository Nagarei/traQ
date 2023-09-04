package handler

import (
	"fmt"
	"github.com/gofrs/uuid"
	"time"

	"github.com/leandro-lugaresi/hub"

	"github.com/traPtitech/traQ/service/bot/event"
	"github.com/traPtitech/traQ/service/bot/event/payload"
)

func UserGroupMemberRemoved(ctx Context, datetime time.Time, _ string, fields hub.Fields) error {
	groupID := fields["group_id"].(uuid.UUID)
	userID := fields["user_id"].(uuid.UUID)
	bots, err := ctx.GetBots(event.UserGroupMemberRemoved)
	if err != nil {
		return fmt.Errorf("failed to GetBots: %w", err)
	}
	if len(bots) == 0 {
		return nil
	}

	if err := ctx.Multicast(
		event.UserGroupMemberRemoved,
		payload.MakeUserGroupMemberRemoved(datetime, groupID, userID),
		bots,
	); err != nil {
		return fmt.Errorf("failed to multicast: %w", err)
	}
	return nil
}
