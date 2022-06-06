package dto

type SignUpRequest struct {
	Userphone string
	Password  string
}

type LoginRequest struct {
	Userphone string
	Password  string
}

type NewDirRequest struct {
	DirName string `json:"dir_name"`
}

type DeleteDirRequest struct {
	DirId int64 `json:"dir_id"`
}

type UpdateDirNameRequest struct {
	NewDirName string `json:"new_dir_name"`
	DirId      int64 `json:"dir_id"`
}

//type GetDirsRequest struct {
//	UserId string `json:"user_id"`
//}

type GetCollectionRequest struct {
	DirId int64 `json:"dir_id"`
}

type CollectRequest struct {
	DirId   int64  `json:"dir_id"`
	Word    string `json:"word"`
	UrlName string `json:"url_name"`
}

type CancelCollectRequest struct {
	CollectionId int64 `json:"collection_id"`
	DirId        int64 `json:"dir_id"`
}
