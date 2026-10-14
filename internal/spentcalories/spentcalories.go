package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)
// parseTraining разбирает строку с данными о тренировке.
// На вход приходит строка вида "3456,Бег,3h00m".
// Возвращает количество шагов, вид активности, длительность и ошибку, если данные некорректны.
func parseTraining(data string) (int, string, time.Duration, error) {
// Разделяю строку на части с помощью запятой
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
        return 0, "", 0, errors.New("неверный формат данных")
	}

    steps, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, "", 0, err
    }
    if steps <= 0 {
    		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
    	}

    duration, err := time.ParseDuration(parts[2])
    if err != nil {
        return 0, "", 0, err
    }
    if duration <= 0 {
    		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
    	}
    return steps, parts[1], duration, nil
}
// distance рассчитывает дистанцию в километрах.
// Принимает количество шагов и рост пользователя, и вычисляет рост.
func distance(steps int, height float64) float64 {
// Вычисляю длину шага
	stepLength := height * stepLengthCoefficient
	    return (float64(steps) * stepLength) / mInKm
}
// meanSpeed рассчитывает среднюю скорость.
// Принимает шаги, рост и длительность.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
// Проверяю что длительность больше нуля
	if duration <= 0 {
        return 0
	}
// Здесь считаем дистанцию и делим на время в часах
    dist := distance(steps, height)
        return dist / duration.Hours()
}
// TrainingInfo формирует строку с полной информацией о тренировке.
// В зависимости от типа (Бег или Ходьба) считает дистанцию, скорость и калории.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
        log.Println(err)
        return "", err
	}
// В зависимости от типа тренировки считаю и оформляю определенную строку
    switch activityType{
    case "Бег":
        calories, err := RunningSpentCalories(steps, weight, height, duration)
        if err != nil {
                return "", err
        }
        dist := distance (steps, height)
        speed := meanSpeed(steps, height, duration)
        return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
         activityType, duration.Hours(), dist, speed, calories), nil

    case "Ходьба":
        calories, err := WalkingSpentCalories(steps, weight, height, duration)
        if err != nil{
            return "", err
        }
        dist := distance(steps, height)
        speed := meanSpeed(steps, height, duration)
        return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
         activityType, duration.Hours(), dist, speed, calories), nil

    default:
        return "", errors.New("неизвестный тип тренировки")
    }
}
// RunningSpentCalories рассчитывает количество калорий, потраченных при беге.
// Возвращает калории и ошибку, если параметры некорректны.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
// Проверка входных параметров на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
        return 0, errors.New("некорректные входные параметры")
	}
// Считаем среднюю скорость
    speed := meanSpeed(steps, height, duration)
// Переводим длительность в минуты
    durationInMinutes := duration.Minutes()

    calories := (weight * speed * durationInMinutes) / minInH

    return calories, nil
}
// WalkingSpentCalories рассчитывает количество калорий, потраченных при ходьбе.
// Работает так же, как и для бега, но с поправочным коэффициентом.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
        return 0, errors.New("некорректные входные параметры")
	}
    speed := meanSpeed(steps, height, duration)

    durationInMinutes := duration.Minutes()

    calories := (weight * speed * durationInMinutes) / minInH
    calories *= walkingCaloriesCoefficient

    return calories, nil
}
