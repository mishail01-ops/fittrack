package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	sc "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	spl := strings.Split(data, ",")
	var step int
	var err error
	var duration time.Duration

	if len(spl) == 2 {

		step, err = strconv.Atoi(spl[0])
		if err != nil {
			return 0, 0, err
		}
		if step <= 0 {
			return 0, 0, fmt.Errorf("Количество шагов должна быть больше нуля")
		}

		duration, err = time.ParseDuration(spl[1])
		if err != nil {
			return 0, 0, err
		}
		if duration <= 0 {
			return 0, 0, fmt.Errorf("Длительность должна быть больше нуля")
		}

		return step, duration, nil
	}

	err = fmt.Errorf("Неверный формат данных")
	return 0, 0, err

}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка DayActionInfo: ", err)
		return ""
	}

	distance := float64(steps) * stepLength / mInKm

	wspent, err := sc.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return ""
	}

	//Количество шагов: 792.
	//Дистанция составила 0.51 км.
	//Вы сожгли 221.33 ккал.
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distance, wspent)

}
