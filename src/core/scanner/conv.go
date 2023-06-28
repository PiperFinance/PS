package scanner

type UserRequest struct {
	Users  []string `json:"users"`
	Tokens []string `json:"tokens"`
	Chains []int64  `json:"chains"`
}
