package main

import (
    "flag"
    "fmt"
    "math/rand"
    "sync"
)

// ProgramParameters - структура параметров программы
type ProgramParameters struct {
    limit int // количество чисел, которые надо сгенерировать
    goNum int // количество горутин, которые пыхтят над генерацией чисел
}

// структура сет - нужно чтобы элементы сравнивались между собой, поэтому comparable
type Set[T comparable] struct {
    data map[T]bool
    lock sync.Mutex
}
//тут есть ресивер
func (s *Set[T]) Add(item T) { //s *Set[T] процедура объявлена на структуре сет.
    s.lock.Lock()
    defer s.lock.Unlock()
    s.data[item] = true // Т.к. мы объявили структуру, мы можем юзать её в типах аргументов
}

func (s *Set[T]) Len() int { // возвращает количество элементов в сете
    s.lock.Lock()
    defer s.lock.Unlock()
    return len(s.data)

}
func (s *Set[T]) ForEach(operation func(item T)) { // проходимся по data(map) и вызываем ф-ю из параметра с каждым элементом отображения
    s.lock.Lock()
    defer s.lock.Unlock()
    for i := range(s.data) {
        operation(i)
    }
}
/*
func (s *Set[T]) get() {
    s.lock.Lock()
    defer s.lock.Unlock()
}*/
// getProgramParameters - получение параметров программы
func getProgramParameters() ProgramParameters {
    limit := flag.Int("limit", 150, "Количество чисел для генерации")
    goNum := flag.Int("goNum", 3, "Количество горутин для генерации чисел")

    flag.Parse()

    return ProgramParameters{limit: *limit, goNum: *goNum}
}

func main() {
    parameters := getProgramParameters()
    numbers := generateNumbers(parameters)
    print(numbers)
}

// generateNumbers -генерирует числа в горутине
func generateNumbers(parameters ProgramParameters) Set[int] {
    numbers := Set[int]{data: make(map[int]bool)}
    //numbers := make(map[int]bool)
    var wg sync.WaitGroup
    limit := parameters.limit
    goNum := parameters.goNum
    for i := 0; i < goNum; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                randomNumber := rand.Intn(limit + 1)
                //lock.Lock() //блокирую мапу чтобы не писали в неё другие потоки, пока читаю
                numbersLen := numbers.Len()
                if numbersLen >= limit {
                    //lock.Unlock()
                    break
                }
                numbers.Add(randomNumber) //если всё ок, то записываю число в отображение
                //lock.Unlock()
            }
        }()
    }
    wg.Wait()
    return numbers
}

// print - печатает сгенерированные числа
func print(numbers Set[int]) {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        numbers.ForEach(func(item int) {   // печать
            fmt.Println(item)
        })
    }()
    wg.Wait()
}