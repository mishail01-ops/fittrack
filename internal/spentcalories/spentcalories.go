package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	//lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	spl := strings.Split(data, ",")

	if len(spl) != 3 {
		return 0, "", 0, fmt.Errorf("Неверный формат данных")
	}
	steps, err := strconv.Atoi(spl[0])
	if err != nil {
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(spl[2])
	if err != nil {
		return 0, "", 0, err
	}

	return steps, spl[1], duration, nil
}

func distance(steps int, height float64) float64 {
	steplen := height * stepLengthCoefficient
	return float64(steps) * steplen / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Hours() <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()

}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, tp, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if weight <= 0 {
		err = fmt.Errorf("Вес должен быть больше нуля")
		log.Println(err)
		return "", err
	}

	if height <= 0 {
		err = fmt.Errorf("Рост должен быть больше нуля")
		log.Println(err)
		return "", err
	}

	var calories float64

	switch tp {
	case "Бег":

		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}

	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", tp)
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	//Тип тренировки: Бег
	//Длительность: 0.75 ч.
	//Дистанция: 10.00 км.
	//Скорость: 13.34 км/ч
	//Сожгли калорий: 18621.75

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", tp, duration.Hours(), dist, speed, calories), nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, fmt.Errorf("Вес должен быть больше нуля")
	}

	if height <= 0 {
		return 0, fmt.Errorf("Рост должен быть больше нуля")
	}

	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов должна быть больше нуля")
	}

	if duration.Minutes() <= 0 {
		return 0, fmt.Errorf("Длительность тренировки должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	calories := (weight * speed * duration.Minutes()) / minInH

	return calories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, fmt.Errorf("Вес должен быть больше нуля")
	}

	if height <= 0 {
		return 0, fmt.Errorf("Рост должен быть больше нуля")
	}

	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов не может быть отрицательным")
	}

	if duration.Minutes() <= 0 {
		return 0, fmt.Errorf("Длительность тренировки должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	calories := walkingCaloriesCoefficient * (weight * speed * duration.Minutes()) / minInH

	return calories, nil
}
