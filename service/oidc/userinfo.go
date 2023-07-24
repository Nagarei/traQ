package oidc

import (
	"github.com/gofrs/uuid"

	"github.com/traPtitech/traQ/model"
	"github.com/traPtitech/traQ/repository"
	"github.com/traPtitech/traQ/service/rbac"
	"github.com/traPtitech/traQ/utils"
)

type Service struct {
	repo   repository.Repository
	origin string
	rbac   rbac.RBAC
}

func NewOIDCService(
	repo repository.Repository,
	origin string,
	rbac rbac.RBAC,
) *Service {
	return &Service{
		repo:   repo,
		origin: origin,
		rbac:   rbac,
	}
}

func (s *Service) GetUserInfo(userID uuid.UUID) (map[string]any, error) {
	user, err := s.repo.GetUser(userID, true)
	if err != nil {
		return nil, err
	}
	tags, err := s.repo.GetUserTagsByUserID(user.GetID())
	if err != nil {
		return nil, err
	}
	groups, err := s.repo.GetUserBelongingGroupIDs(user.GetID())
	if err != nil {
		return nil, err
	}

	return map[string]any{
		// OIDC standard claims
		"name":               user.GetName(),
		"preferred_username": user.GetName(),
		"picture":            s.origin + "/api/v3/public/icon/" + user.GetName(),
		"updated_at":         user.GetUpdatedAt(),
		// traQ specific claims
		"traq": map[string]any{
			"bio":          user.GetBio(),
			"groups":       groups,
			"tags":         utils.Map(tags, func(tag model.UserTag) string { return tag.GetTag() }),
			"last_online":  user.GetLastOnline(),
			"twitter_id":   user.GetTwitterID(),
			"display_name": user.GetResponseDisplayName(),
			"icon_file_id": user.GetIconFileID(),
			"bot":          user.IsBot(),
			"state":        user.GetState().Int(),
			"permissions":  s.rbac.GetGrantedPermissions(user.GetRole()),
			"home_channel": user.GetHomeChannel(),
		},
	}, nil
}
