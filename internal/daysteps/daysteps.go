package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/sprint4Final/tracker/internal/spentcalories"
	"log"


)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage принимает строку, возвращает количество шагов и продолжительность прогулки
func parsePackage(data string) (int, time.Duration, error) {
	separator := ","
	sliceData := strings.Split(data, separator)
	if len(sliceData) != 2 {
		fmt.Println("Неверная длина слайса")
		return 0, 0, fmt.Errorf("неверная длина данных")
	}
	countSteps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		return 0, 0, err
	}
	if countSteps <= 0 {
		fmt.Println("Неверное количество шагов")
		return 0, 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}
	walkDuration, err := time.ParseDuration(sliceData[1])
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		return 0, 0, err
	}
	if walkDuration <= 0 {
		fmt.Println("Неверное время")
		return 0, 0, fmt.Errorf("неверное время")
	}
	return countSteps, walkDuration, nil
}

// DayActionInfo вычисляет дистанцию в километрах и количество потраченных калорий
func DayActionInfo(data string, weight, height float64) string {
    countSteps, walkDuration, err := parsePackage(data)
    if err != nil {
        log.Printf("Ошибка при парсинге данных: %v\n", err)
        return ""
    }
    if countSteps <= 0 {
        log.Println("Неверное количество шагов")
        return ""
    }
    distance := stepLength * float64(countSteps)
    distance = distance / float64(mInKm)
    calories, err := spentcalories.WalkingSpentCalories(countSteps, weight, height, walkDuration)
    if err != nil {
        log.Printf("Ошибка при расчёте сожжённых калорий: %v\n", err)
        return ""
    }
    return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", countSteps, distance, calories)
}

