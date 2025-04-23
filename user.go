package main

type userRequest struct {
	Name string `json:"name"`
}

type user struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
}
