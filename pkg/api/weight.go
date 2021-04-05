package api

import (
	"errors"
	"time"
)

type WeightService interface {
	New(request NewWeightRequest) error
	CalculateBMR(height, age, weight int, sex string) (int, error)
	DailyIntake(BMR, activityLevel int, weightGoal string) (int, error)
}

type WeightRepository interface {
	CreateWeightEntry(w Weight) error
	GetUser(userID int) (User, error)
}

type weightService struct {
	storage WeightRepository
}

func NewWeightService(weightRepo WeightRepository) WeightService {
	return &weightService{storage: weightRepo}
}

type User struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Age           int       `json:"age"`
	Height        int       `json:"height"`
	Sex           string    `json:"sex"`
	ActivityLevel int       `json:"activity_level"`
	WeightGoal    string    `json:"weight_goal"`
	Email         string    `json:"email"`
}

type Weight struct {
	Weight             int `json:"weight"`
	UserID             int `json:"user_id"`
	BMR                int `json:"bmr"`
	DailyCaloricIntake int `json:"daily_caloric_intake"`
}

type NewWeightRequest struct {
	Weight int `json:"weight"`
	UserID int `json:"user_id"`
}

const (
	// Very Low Intensity activity multiplier - 1
	veryLowActivity = 1.2
	// Light exercise activity multiplier - 3-4x per week - 2
	lightActivity = 1.375
	// Moderate exercise activity multiplier - 3-5x per week 30-60 mins/session - 3
	moderateActivity = 1.55
	// High exercise activity multiplier - (6-7x per week for 45-60 mins) - 4
	highActivity = 1.725
	// Very high exercise activity multiplier - for people who is an Athlete - 5
	veryHighActivity = 1.9
)

func (w *weightService) New(request NewWeightRequest) error {
	if request.UserID == 0 {
		return errors.New("weight service - user ID cannot be 0")
	}

	user, err := w.storage.GetUser(request.UserID)

	if err != nil {
		return err
	}

	bmr, err := w.CalculateBMR(user.Height, user.Age, request.Weight, user.Sex)

	if err != nil {
		return err
	}

	dailyIntake, err := w.DailyIntake(bmr, user.ActivityLevel, user.WeightGoal)

	if err != nil {
		return err
	}

	newWeight := Weight{
		Weight:             request.Weight,
		UserID:             user.ID,
		BMR:                bmr,
		DailyCaloricIntake: dailyIntake,
	}

	err = w.storage.CreateWeightEntry(newWeight)

	if err != nil {
		return err
	}

	return nil
}

func (w *weightService) CalculateBMR(height, age, weight int, sex string) (int, error) {
	var sexModifier int

	switch sex {
	case "male":
		sexModifier = -5
	case "female":
		sexModifier = 161
	default:
		return 0, errors.New("invalid variable sex provided to CalculateBMR. needs to be either male or female")
	}

	return (10 * weight) + int(float64(height)*6.25) - (5 * age) - sexModifier, nil
}

func (w *weightService) DailyIntake(BMR, activityLevel int, weightGoal string) (int, error) {
	var maintenanceCalories int

	switch activityLevel {
	case 1:
		maintenanceCalories = int(float64(BMR) * veryLowActivity)
	case 2:
		maintenanceCalories = int(float64(BMR) * lightActivity)
	case 3:
		maintenanceCalories = int(float64(BMR) * moderateActivity)
	case 4:
		maintenanceCalories = int(float64(BMR) * highActivity)
	case 5:
		maintenanceCalories = int(float64(BMR) * veryHighActivity)
	default:
		return 0, errors.New("invalid variable activityLevel - needs to be 1, 2, 3, 4 or 5")
	}

	var dailyCaloricIntake int

	switch weightGoal {
	case "gain":
		dailyCaloricIntake = maintenanceCalories + 500
	case "loose":
		dailyCaloricIntake = maintenanceCalories - 500
	case "maintain":
		dailyCaloricIntake = maintenanceCalories
	default:
		return 0, errors.New("invalid weight goal provided - must be gain, loose or maintain")
	}

	return dailyCaloricIntake, nil
}
