package dto

import (
	"gofound/common"
	"gofound/models"
)

type CreateUserResponse struct {
	Code common.ErrNo
	Data struct {
		UserID string // int64 范围
	}
}
type SignUpResponse struct {
	Code common.ErrNo
	Data struct {
		UserID string
	}
}

type LoginResponse struct {
	Code common.ErrNo
	Data struct {
		UserID string
	}
}

type LogoutResponse struct {
	Code common.ErrNo
}

type WhoAmIResponse struct {
	Code common.ErrNo
	Data models.TUser
}

type DeleteUserResponse struct {
	Code common.ErrNo
}

type NewDirResponse struct {
	Code common.ErrNo
	Data struct {
		DirId string
	}
}

type DeleteDirResponse struct {
	Code common.ErrNo
}

type UpdateDirNameResponse struct {
	Code common.ErrNo
}

type GetDirResponse struct {
	Code common.ErrNo
	Data struct {
		DirList []models.Dir
	}
}
type GetCollectionResponse struct {
	Code common.ErrNo
	Data struct {
		CollectionList []models.Collection
	}
}
type CollectResponse struct {
	Code common.ErrNo
	Data struct {
		CollectionID string
	}
}

type CancelCollectResponse struct {
	Code common.ErrNo
}
