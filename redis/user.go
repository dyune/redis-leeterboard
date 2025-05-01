package redis

type UserRequest struct {
	Name string `json:"name"`
}

type User struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Rank  int     `json:"rank"`
	Score float64 `json:"score"`
}
