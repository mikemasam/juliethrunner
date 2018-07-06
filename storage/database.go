package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"juliethrunner/dateutils"
	"juliethrunner/utils"
	"sort"
	"time"
)

type Database struct {
	DB_path     string
	Config_path string
	Configs     []Configs
	Tasks       []Tasks
}

//var DB_path string = "data/tasks.db"

func (d *Database) InitTasks() {

	_tasks := []Tasks{}
	fp := utils.OpenFile(d.DB_path)
	jsonParser := json.NewDecoder(fp)
	if err := jsonParser.Decode(&_tasks); err != nil {
		fmt.Println("Data not available")
	}

	fp.Close()

	d.Tasks = _tasks
}

func (d *Database) InitConfigs() {

	_configs := []Configs{}
	fp := utils.OpenFile(d.Config_path)
	jsonParser := json.NewDecoder(fp)
	if err := jsonParser.Decode(&_configs); err != nil {
		fmt.Println("parsing config file", err.Error())
	}

	fp.Close()

	d.Configs = _configs
}

func (d *Database) SaveTasks() {

	//	for _, c := range d.Configs {
	//		fmt.Println(c)
	//	}

	for i := range d.Tasks {
		d.Tasks[i].Enabled = false
	}

	for _, v := range d.Configs {
		_task := Tasks{}
		_task.Command = v.Commands
		_task.Id = v.Id
		_task.Last_run = time.Now().Unix()
		_task.Interval, _task.Repeateble = dateutils.GetUnixTimeFromZero(v.Every)

		d.addTask(_task)
	}

	b, _ := json.Marshal(d.Tasks)

	fp := utils.OpenFile(d.DB_path)
	fp.WriteString(string(b))
	fp.Close()
	//fmt.Println(string(b))
}

func (d *Database) PopulateCron(app_name string) {
	fp := utils.OpenFileWithPermission(app_name, 0644)
	if fp != nil {

		for _, v := range d.Tasks {
			fp.WriteString("*/1 * * * * root " + app_name + " -id " + v.Id + "\n")
		}
		fp.Close()
	}
}

func (d *Database) addTask(task Tasks) {
	if !d.TaskExists(task.Id) {
		task.Enabled = true
		d.Tasks = append(d.Tasks, task)
	}
}
func (d *Database) removeTaskByIndex(i int) {

	a := d.Tasks

	a[i] = a[len(a)-1] // Copy last element to index i
	a = a[:len(a)-1]   // Truncate slice
	d.Tasks = a
}
func (t *Database) TaskExists(id string) bool {
	for i := range t.Tasks {
		if t.Tasks[i].Id == id {
			if !t.Tasks[i].Enabled {
				t.removeTaskByIndex(i)
				return false
			}
			return true
		}
	}
	return false
}

func (t *Database) FindTaskById(id string) (Tasks, error) {
	for _, v := range t.Tasks {
		if v.Id == id {
			return v, nil
		}
	}
	return Tasks{}, errors.New("Task not found")
}

func (t *Database) findAndRunOlder(_time int64) bool {

	got_one := false

	if len(t.Tasks) > 0 {
		v := &t.Tasks[0]
		if v.When() <= _time {
			v.Called()

			go func() {
				t.SaveTasks()
				sort.Sort(TaskSortable(t.Tasks))
				t.RunTask(*v)
			}()

			got_one = true
		}
	}

	return got_one
}

func (d *Database) RunTask(t Tasks) {
	utils.RunTasks(t.Command)
}

func (d *Database) Run() {

	sort.Sort(TaskSortable(d.Tasks))

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		current := t.Unix()

		d.findAndRunOlder(current)

		fmt.Println("Tick at", current)
	}
}

func (d *Database) hasTask() bool {
	return len(d.Tasks) > 0
}
