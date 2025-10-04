package spentcalories

import (
	"time"
	"fmt"
	"strconv"
	"strings"
	"log"
	"errors"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining принимает строку с данными и возвращает количество шагов, тип активности, время активности
func parseTraining(data string) (int, string, time.Duration, error) {
	dataWithoutSpaces := strings.Replace(data, " ", "", -1)
	sliceData := strings.Split(dataWithoutSpaces, ",")
	if len(sliceData) != 3 {
		fmt.Println("Неверная длина слайса")
		return 0, "", 0, fmt.Errorf("неверная длина данных")
	}
	countSteps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		return 0, "", 0, err
	}
	if countSteps <= 0 {
		fmt.Println("Неверное количество шагов")
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}
	activityDuration, err := time.ParseDuration(sliceData[2])
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		return 0, "", 0, err
	}
	if activityDuration <= 0 {
    fmt.Println("Неверное время")
    return 0, "", 0, fmt.Errorf("продолжительность не может быть отрицательной или нулевой")
    }
	typeActivity := sliceData[1]
	return countSteps, typeActivity, activityDuration, nil
}

// distance принимает количество шагов и рост пользователя в метрах, а возвращает дистанцию в километрах
func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	dist := stepLen * float64(steps)
	dist = dist / mInKm
	return dist
}

// meanSpeed принимает количество шагов, рост пользователя и продолжительность активности и возвращает среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / float64(duration.Hours())
}

// TrainingInfo возвращает строку с информацией о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	countSteps, typeActivity, activityDuration, err := parseTraining(data)
    if err != nil {
        fmt.Println("Ошибка:", err)
		log.Println(err)
        return "", err
	}
	switch typeActivity {
	case "Ходьба":
		dist := distance(countSteps,height)
		midSpeed := meanSpeed(countSteps, height, activityDuration)
		calories, err := WalkingSpentCalories(countSteps, weight, height, activityDuration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
        typeActivity, float64(activityDuration.Hours()), dist, midSpeed, calories), nil
	case "Бег":
		dist := distance(countSteps,height)
		midSpeed := meanSpeed(countSteps, height, activityDuration)
		calories, err := RunningSpentCalories(countSteps, weight, height, activityDuration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
        typeActivity, float64(activityDuration.Hours()), dist, midSpeed, calories), nil
	default:
        return "", fmt.Errorf("неизвестный тип тренировки: %s", typeActivity)
	}
	
}

// RunningSpentCalories возвращает количество калорий, потраченных при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
        return 0, errors.New("все входные параметры должны быть больше нуля")
    }
	midSpeed := meanSpeed(steps, height, duration)
	minutesDuration := duration.Minutes()
	return (weight * midSpeed * minutesDuration) / minInH, nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при беге...
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
        return 0, errors.New("все входные параметры должны быть больше нуля")
    }
	midSpeed := meanSpeed(steps, height, duration)
	minutesDuration := duration.Minutes()
	return (weight * midSpeed * minutesDuration * walkingCaloriesCoefficient) / minInH, nil
}


