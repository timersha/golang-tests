package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrCounter struct {
	mx    sync.RWMutex
	value int
}

func (ec *ErrCounter) Increment() {
	ec.mx.Lock()
	defer ec.mx.Unlock()
	ec.value++
}

func (ec *ErrCounter) Read() int {
	ec.mx.RLock()
	defer ec.mx.RUnlock()
	return ec.value
}

func Run(tasks []Task, n, m int) error {
	done := make(chan struct{}) // канал через который задача сообщает о своем завершении
	defer close(done)

	var wg sync.WaitGroup       // отслеживание завершения задач, при превышении количества ошибок
	var errResult error         // возварщаемыйрезультат
	errCount := ErrCounter{}    // количество ошибок
	batch := tasks[:n]          // первая пачка задач для выполнения
	startCount := n             // количество запущенных задач
	doneCount := 0              // количество завершенных задач
	allTasksCount := len(tasks) // количество всех задач
	m = max(m, 1)

	// Обертка для запуска задачи
	task := func(done chan<- struct{}, t Task, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		e := t()
		if e != nil {
			// Одновременно несколько горутин может изменять значение errCount
			// поэтому изменение значения errCount происходит через RWMutex
			errCount.Increment()
		}
		done <- struct{}{}
	}
	// Запускаем первую пачку задач
	for _, t := range batch {
		go task(done, t, &wg)
	}

	// Нужно прочитать все возможные сообщения от задач
	for i := 0; i < allTasksCount; i++ {
		select {
		case <-done:
			doneCount++

			// Если количество ошибок превысило лимит, то прерываем чтение.
			// Пока читаем errCount уже запущенные горутины могут изменять это значение,
			// поэтому читаем через RWMutex
			if errCount.Read() >= m {
				errResult = ErrErrorsLimitExceeded
				// "Дочитываем" уже запущенные задачи и останавливаем цикл чтения
				// иначе горутины будут писать в закрытый канал
				allTasksCount = startCount
				break
			}

			// Все задачи сообщили о своем завершении
			if doneCount == len(tasks) {
				return nil
			}

			// Как только одна задача завершилась, запускаем новую пока есть возможность
			if startCount < len(tasks) && errCount.Read() < m {
				go task(done, tasks[startCount], &wg)
				startCount++
			}
		}
	}
	wg.Wait()
	return errResult
}
