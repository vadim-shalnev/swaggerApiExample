package Models

type RequestSearch struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type RequestAddress struct {
	Addres        string        `json:"addres"`
	RequestSearch RequestSearch `json:"search"`
}
type RequestQuery struct {
	Query string `json:"query"`
}
type RequestUser struct {
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	RequestQuery RequestQuery `json:"request_query"`
}
type SearchHistory struct {
	ID      int    `json:"id"`
	Query   string `json:"query"`
	Address string `json:"address"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
}
