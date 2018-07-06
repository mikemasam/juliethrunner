package storage

import (
	"fmt"
	"time"
)

type Tasks struct {
	Last_run   int64    `json:"last_run"`
	Interval   int64    `json:"interval"`
	Command    []string `json:"command"`
	Id         string   `json:"id"`
	Repeat     bool     `json:"repeat"`
	Repeateble bool     `json:"repeateble"`
	Enabled    bool     `json:"enabled"`
}

func (d *Tasks) When() int64 {
	return d.Last_run + d.Interval
}
func (d *Tasks) Called() {
	d.Last_run = time.Now().Unix()
	fmt.Println(d.Id+" called on ", time.Now())
}

type TaskSortable []Tasks

func (a TaskSortable) Len() int           { return len(a) }
func (a TaskSortable) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TaskSortable) Less(i, j int) bool { return a[i].When() < a[j].When() }
