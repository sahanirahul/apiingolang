package irepo

type IActivityRepo interface {
	InsertActivities() error
}

type IHttpRepo interface {
	CallExternal() error
}
