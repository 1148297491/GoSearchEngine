package common

type ErrNo int

const (
	OK                   ErrNo = 0
	ParamInvalid         ErrNo = 1   // 参数不合法
	UserHasExisted       ErrNo = 2   // 该 Userphone 已存在
	UserHasDeleted       ErrNo = 3   // 用户已删除
	UserNotExisted       ErrNo = 4   // 用户不存在
	WrongPassword        ErrNo = 5   // 密码错误
	LoginRequired        ErrNo = 6   // 用户未登录
	PermDenied           ErrNo = 10  // 没有操作权限
	DirHasExisted        ErrNo = 11  // 收藏夹已存在
	DirNotExisted        ErrNo = 12  // 收藏夹不存在
	CollectionHasExisted ErrNo = 13  //该收藏夹中已存在该搜索结果
	CollectionNotExisted ErrNo = 14  //该结果不存在
	UserHaveNoDir        ErrNo = 15  //用户收藏夹为空
	UnknownError         ErrNo = 255 // 未知错误
)
