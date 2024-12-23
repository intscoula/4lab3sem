package main

import (
	"fmt"
	"math/rand"
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
	if avgAge > 0 {
		fmt.Printf("Средний возраст женщин: %.2f лет\n", avgAge)
	} else {
		fmt.Printf("Нет женщин, родивших в указанный период.\n")
	}
	fmt.Printf("Время обработки: %v\n\n", duration)
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
		fmt.Printf("Нет женщин, родивших в указанный период.\n")
		fmt.Printf("Время обработки: %v\n\n", time.Since(start))
		return
	}

	avgAge := totalAge.Hours() / 24 / 365.25 / float64(count)

	duration := time.Since(start)

	fmt.Printf("С многозадачностью:\n")
	if avgAge > 0 {
		fmt.Printf("Средний возраст женщин: %.2f лет\n", avgAge)
	} else {
		fmt.Printf("Нет женщин, родивших в указанный период.\n")
	}
	fmt.Printf("Время обработки: %v\n\n", duration)
}

// Генерация случайных данных о роддоме
func generateMaternity(index int, startDate, endDate time.Time) Maternity {
	// Генерируем имя по шаблону
	name := fmt.Sprintf("Женщина %d", index)
	// Генерируем случайную дату рождения женщины (от 18 до 45 лет назад)
	womanBirthDate := time.Now().AddDate(-rand.Intn(28)-18, -rand.Intn(12), -rand.Intn(31))
	// Генерируем случайную дату рождения ребенка в указанный период
	childBirthDate := randomDateInRange(startDate, endDate)

	return Maternity{
		WomanName:      name,
		WomanBirthDate: womanBirthDate,
		ChildBirthDate: childBirthDate,
	}
}

// Генерация случайной даты в указанном диапазоне
func randomDateInRange(startDate, endDate time.Time) time.Time {
	diff := endDate.Sub(startDate)
	randomDuration := time.Duration(rand.Int63n(int64(diff)))
	return startDate.Add(randomDuration)
}

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Период для обработки
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// Массив данных о роддоме размером 100000
	var maternities []Maternity
	for i := 0; i < 100000; i++ {
		maternities = append(maternities, generateMaternity(i, startDate, endDate))
	}

	// Обработка без многозадачности
	processWithoutConcurrency(maternities, startDate, endDate)

	// Обработка с многозадачностью
	processWithConcurrency(maternities, startDate, endDate)
}
