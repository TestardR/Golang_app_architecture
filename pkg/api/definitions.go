package api

type NewUserRequest struct {
	Name          string `json:"name"`
	Age           int    `json:"age"`
	Height        int    `json:"height"`
	Sex           string `json:"sex"`
	ActivityLevel int    `json:"activity_level"`
	WeightGoal    string `json:"weight_goal"`
	Email         string `json:"email"`
}
