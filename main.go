package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// WaitGroup - Наша простенькая реализация объекта ожидания завершения группы горутин
type WaitGroup struct {
	counter int
	// Условная переменная, через которую сам
	// объект будет уведомлять другие горутины,
	// что нужные им горутины завершили работу.
	// Используется внутренний локер для потокобезопасного доступа
	// к полю counter,
	// так как counter - разделяемый между горутинами ресурс
	c *sync.Cond
}

// NewWaitGroup - создание объекта ожидания завершения группы горутин
func NewWaitGroup() *WaitGroup {
	return &WaitGroup{0, sync.NewCond(&sync.Mutex{})}
}

// Add - добавление определенного количества
// горутин для ожидания
// завершения
func (w *WaitGroup) Add(amount int) {
	w.c.L.Lock()
	w.counter += amount
	w.c.L.Unlock()
}

// Done - вызов данного метода сигнализирует о завершении работы одной из
// горутин
func (w *WaitGroup) Done() {
	w.c.L.Lock()
	w.counter--
	w.c.L.Unlock()
	w.c.Broadcast()
}

// Wait -  ожидание завершения работы добавленного количества горутин
func (w *WaitGroup) Wait() {
	defer w.c.L.Unlock()
	w.c.L.Lock()
	for w.counter != 0 {
		w.c.Wait()
	}
}

// Шаг наращивания счётчика
const step int64 = 1

// Конечное значение счетчика
const endCounterValue int64 = 1000

func main() {

	var counter int64 = 0
	var wg = NewWaitGroup()
	var iterationCount = int(endCounterValue / step)
	increment := func() {
		defer wg.Done()
		for i := 1; i <= iterationCount; i++ {
			atomic.AddInt64(&counter, step)
		}
	}

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go increment()
	}

	wg.Wait()
	fmt.Println(counter)
}
