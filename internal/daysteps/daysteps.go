package daysteps

import (
    "fmt"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)
// parsePackage разбирает строку с данными о прогулке.
// На вход приходит строка вида "678,0h50m".
// Возвращает количество шагов, длительность прогулки и ошибку,если данные некорректны.
func parsePackage(data string) (int, time.Duration, error) {
// Разделяем строку на части с помощью запятой
	parts := strings.Split(data, ",")
// Проверяю, что получится именно две части
	if len(parts) != 2{
	    return 0, 0, errors.New("неверный формат данных")
	}
    steps, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, 0, err
    }
// Проверяю, что количество шагов больше нуля
    if steps <= 0{
        return 0, 0, errors.New("количество шагов должно быть больше 0")
    }
    duration, err := time.ParseDuration(parts[1])
    if err != nil{
        return 0, 0, err
    }
    if duration <= 0 {
    		return 0, 0, errors.New("продолжительность должна быть больше 0")
    	}
    return steps, duration, nil
}
// DayActionInfo формирует итоговую строку с информацией о дневной активности.
// Принимает строку с данными о прогулке, вес и рост пользователя.
// Считает дистанцию и количество сожжённых калорий, и возвращает строку с результатами.
// Возвращает пустую строку, если данные некорректны.
func DayActionInfo(data string, weight, height float64) string {
// Получаю шаги и длительность через parsePackage
	steps, duration, err := parsePackage(data)
	if err != nil {
	return ""
	}
    if steps <= 0{
        return ""
    }
    distance := float64(steps) * stepLength / mInKm

    calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
    if err != nil{
        return ""
    }

    result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
    return result
}
