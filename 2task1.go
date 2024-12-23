package main

import (
	"fmt"
	"sync"
	"time"
)

// Структура содержит сведения о роддоме
type Maternity struct {
	WomanName      string
	WomanBirthDate time.Time
	ChildBirthDate time.Time
}

// Функция для вычисления среднего возраста женщин
func calculateAverageWomanAge(maternities []Maternity, startDate, endDate time.Time) float64 {
	var totalAge time.Duration
	var count int

	for _, maternity := range maternities {
		if maternity.ChildBirthDate.After(startDate) && maternity.ChildBirthDate.Before(endDate) {
			age := maternity.ChildBirthDate.Sub(maternity.WomanBirthDate)
			totalAge += age
			count++
		}
	}

	if count == 0 {
		return 0
	}

	averageAge := totalAge / time.Duration(count)
	return averageAge.Hours() / 24 / 365.25 // Перевод в годы
}

// Обработка без многозадачности
func processWithoutConcurrency(maternities []Maternity, startDate, endDate time.Time) {
	start := time.Now()

	avgAge := calculateAverageWomanAge(maternities, startDate, endDate)

	duration := time.Since(start)

	fmt.Printf("Без многозадачности:\n")
	fmt.Printf("Средний возраст женщин: %.2f лет\n", avgAge)

	// Вывод времени в микросекундах или наносекундах
	if duration.Microseconds() > 0 {
		fmt.Printf("Время обработки: %d микросекунд\n\n", duration.Microseconds())
	} else {
		fmt.Printf("Время обработки: %d наносекунд\n\n", duration.Nanoseconds())
	}
}

// Обработка с многозадачностью
func processWithConcurrency(maternities []Maternity, startDate, endDate time.Time) {
	start := time.Now()

	var wg sync.WaitGroup
	var totalAge time.Duration
	var count int

	// Разбиваем данные на части для многозадачности
	numGoroutines := 3 // Количество горутин
	chunkSize := len(maternities) / numGoroutines

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			startIndex := i * chunkSize
			endIndex := (i + 1) * chunkSize
			if i == numGoroutines-1 { // Для последнего сегмента
				endIndex = len(maternities)
			}
			subset := maternities[startIndex:endIndex]
			subsetAge := calculateAverageWomanAge(subset, startDate, endDate)
			totalAge += time.Duration(subsetAge * 24 * 365.25 * float64(time.Hour))
			count += len(subset)
		}(i)
	}

	wg.Wait()

	if count == 0 {
		fmt.Printf("С многозадачностью:\n")
		fmt.Printf("Средний возраст женщин: 0 лет\n")
		duration := time.Since(start)
		if duration.Microseconds() > 0 {
			fmt.Printf("Время обработки: %d микросекунд\n\n", duration.Microseconds())
		} else {
			fmt.Printf("Время обработки: %d наносекунд\n\n", duration.Nanoseconds())
		}
		return
	}

	avgAge := totalAge.Hours() / 24 / 365.25 / float64(count)

	duration := time.Since(start)

	fmt.Printf("С многозадачностью:\n")
	fmt.Printf("Средний возраст женщин: %.2f лет\n", avgAge)

	// Вывод времени в микросекундах или наносекундах
	if duration.Microseconds() > 0 {
		fmt.Printf("Время обработки: %d микросекунд\n\n", duration.Microseconds())
	} else {
		fmt.Printf("Время обработки: %d наносекунд\n\n", duration.Nanoseconds())
	}
}

func main() {
	// Пример данных
	maternities := []Maternity{
		{"Иванова Анна", time.Date(1985, 5, 10, 0, 0, 0, 0, time.UTC), time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)},
		{"Петрова Мария", time.Date(1990, 8, 20, 0, 0, 0, 0, time.UTC), time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC)},
		{"Сидорова Елена", time.Date(1988, 3, 15, 0, 0, 0, 0, time.UTC), time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC)},
		{"Кузнецова Ольга", time.Date(1992, 11, 5, 0, 0, 0, 0, time.UTC), time.Date(2023, 4, 5, 0, 0, 0, 0, time.UTC)},
		{"Морозова Татьяна", time.Date(1987, 7, 25, 0, 0, 0, 0, time.UTC), time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)},
	}

	// Период для обработки
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)

	// Обработка без многозадачности
	processWithoutConcurrency(maternities, startDate, endDate)

	// Обработка с многозадачностью
	processWithConcurrency(maternities, startDate, endDate)
}
