package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"task_tracker/pkg/storage"
	"task_tracker/pkg/storage/postgres"
)

var db storage.Interface

func main() {
	var err error
	pwd := os.Getenv("cyber")
	host := os.Getenv("ipaddr")
	if pwd == "" {
		os.Exit(1)
	}
	connstr := "postgres://cyber:" + pwd + "@" + host + "/tasktracker"
	db, err = postgres.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	//db = memdb.DB{}
	for i:= 1; i < 11 ; i ++ {
		taskId := strconv.Itoa(i)
		id, err := db.NewTask(postgres.Task{Title: "task #" +
			taskId, Content: "task content #" + taskId, AuthorID: 1, AssignedID: 1})
		if err != nil	 {
			log.Fatal(err)
		}
		fmt.Println("Добавлена задача с id:", id)
	}

	tasks, err := db.Tasks(0, 0)
	if err != nil	 {
		log.Fatal(err)
	}
	fmt.Println(tasks)

	authName := "Doe"
	labelName := "Go"
	taskId := 2
	newContent := "10-12-21"
	deleteTaskId := 3

	fmt.Println("Задачи по имени автора: ", authName)
	authorTasks, err := db.AuthorTasks(authName)
	if err != nil	 {
		log.Fatal(err)
	}
	fmt.Println(authorTasks)

	fmt.Println("Задачи по тегу: ", labelName)
	labelTasks, err := db.LabelTasks(labelName)
	if err != nil	 {
		log.Fatal(err)
	}
	fmt.Println(labelTasks)

	fmt.Println("Обновление текста задачи с id = ", taskId)
	updatedTasks, err := db.UpdateTasks(taskId, newContent)
	if err != nil	 {
		log.Fatal(err)
	}
	fmt.Println(updatedTasks)

	fmt.Println("Удаление задачи с id = ", deleteTaskId)
	deletedTasks, err := db.DeleteTasks(deleteTaskId, )
	if err != nil	 {
		log.Fatal(err)
	}
	fmt.Println(deletedTasks)
}


