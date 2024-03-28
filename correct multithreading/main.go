package main

import (
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life

type Ttype struct {
	id         int
	cT         string // creation time
	fT         string // finish time
	taskResult string
}

func main() {

	taskChan := make(chan Ttype, 10)
	results := make(chan Ttype)
	errors := make(chan error)

	doneTasks := make(map[int]Ttype)
	undoneTasks := make(map[int]Ttype)

	// Мьютекс для синхронизации доступа к os.Stdout
	var mutex sync.Mutex

	// Горутина для обработки результатов и ошибок
	go func() {
		var wg sync.WaitGroup
		for result := range results {
			doneTasks[result.id] = result
			wg.Add(1)
			go func(task Ttype) {
				defer wg.Done()
				mutex.Lock()
				defer mutex.Unlock()
				fmt.Println("Done task:", task.id)
			}(result)
		}
		for err := range errors {
			undoneTasks[err.(*taskError).task.id] = err.(*taskError).task
			mutex.Lock()
			defer mutex.Unlock()
			fmt.Println("Error:", err)
			time.Sleep(time.Millisecond * 100)
		}
		wg.Wait()
	}()

	// Функция для создания заданий
	go func() {
		defer close(taskChan)
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occurred"
			}
			taskChan <- Ttype{cT: ft, id: int(time.Now().Unix())}
		}
	}()

	// Запуск горутин обработки заданий
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				taskWorker(task, results, errors)
			}
		}()
	}
	wg.Wait()

	// Вывод результатов
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println("Done tasks:")
	for id := range doneTasks {
		fmt.Println(id)
	}
	fmt.Println("Undone tasks:")
	for id := range undoneTasks {
		fmt.Println(id)
	}
}

// Функция для обработки заданий
func taskWorker(task Ttype, results chan<- Ttype, errors chan<- error) {
	tt, _ := time.Parse(time.RFC3339, task.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.taskResult = "task has been succeeded"
	} else {
		task.taskResult = "something went wrong"
		errors <- &taskError{task: task, err: fmt.Errorf("Task id %d, creation time %s, error: %s", task.id, task.cT, task.taskResult)}
	}
	task.fT = time.Now().Format(time.RFC3339Nano)
	results <- task
}

// Ошибка задания
type taskError struct {
	task Ttype
	err  error
}

func (te *taskError) Error() string {
	return te.err.Error()
}
