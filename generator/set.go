package main

import (
	"sync"
)

// struct Set - структура, реализующая сет
// Нужно чтобы элементы сравнивались между собой, поэтому comparable
type Set[T comparable] struct {
	data map[T]bool
	lock sync.Mutex
}

// Add - Функция добавления элемента в сет
func (s *Set[T]) Add(item T) bool { //s *Set[T] процедура объявлена на структуре сет.
	s.lock.Lock()
	defer s.lock.Unlock()
	hasItem := s.data[item]
	s.data[item] = true // Т.к. мы объявили структуру, мы можем юзать её в типах аргументов
	return !hasItem
}

// Len() - Возвращает количество элементов в сете
func (s *Set[T]) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.data)

}

// ForEach - Итератор
func (s *Set[T]) ForEach(operation func(item T)) { // проходимся по data(map) и вызываем ф-ю из параметра с каждым элементом отображения
	s.lock.Lock()
	defer s.lock.Unlock()
	for i := range s.data {
		operation(i)
	}
}

