package main

import (
	"errors"
	"flag"
	"fmt"
	"juliethrunner/storage"
	"juliethrunner/utils"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
)

var app_name string = "JuliethRunner"
var (
	signal        = flag.String("s", "", `send signal to the daemon stop,reload`)
	stop          bool
	compile       = flag.Bool("compile", false, "compile config file to database")
	list_tasks    = flag.Bool("list-tasks", false, "list tasks from database")
	config_file   = flag.String("conf", "~/.julieth.json", "Configurations file")
	database_file = flag.String("db", "~/.julieth.db", "Configurations Database")
	task_id       = flag.String("id", "", "Task ID to run")
	_service      = flag.Bool("service", false, "Running daemon")
	debug         = flag.Bool("debug", false, "Debuging mode")
	database      = storage.Database{}
)

func main() {

	pid := os.Getpid()
	fmt.Println("welcome to julieth runner PiID = ", pid)
	flag.BoolVar(&stop, "stop", false, "Terminate daemon")
	flag.Parse()

	dmn := handleDaemonSignals()

	//fmt.Println(len(os.Args), os.Args)

	handleCompile()
	handleDatabaseQueries()

	//	if *run {
	//		//er := utils.RunDaemon([]string{os.Args[0], "-service true"})
	//	}

	if *_service && dmn != nil {

		// Process daemon operations - send signal if present flag or daemonize
		child, err := dmn.Reborn()
		if err != nil {
			fmt.Println("Start failed = ", err)
			log.Fatalln(err)
		}
		if child != nil {
			fmt.Println(errors.New("Daemon Child Active"))
			return
		}
		defer dmn.Release()

		fmt.Println("Service started")

		work()
		go func() {
			//			for {
			//				time.Sleep(time.Second)
			//				fmt.Println("Running = ", time.Now())
			//			}
		}()

		err = daemon.ServeSignals()
		if err != nil {
			log.Println("Error:", err)
		}

		//		er := utils.RunDaemon([]string{"-id 123 -service=true"})
		//		fmt.Println("Daemon Response := ", er)

	}
}

var working bool = false

func work() {
	go _work()
}
func _work() {

	if working {
		return
	}

	working = true
	fmt.Println("Daemon is running")
	loadDatabase()
	if len(database.Tasks) > 0 {
		database.Run()
	} else {
		fmt.Println("Empty tasks")
	}
	working = false
}

func loadDatabase() {
	database = storage.Database{}
	database.Config_path = *config_file
	database.DB_path = *database_file
	database.InitTasks()
}
func SysHandler(sig os.Signal) error {
	fmt.Println("signal:", sig)
	if sig == syscall.SIGTERM {
		return daemon.ErrStop
	}

	if sig == syscall.SIGHUP {
		work()
	}

	return nil
}

func handleDaemonSignals() *daemon.Context {

	daemon.AddCommand(daemon.BoolFlag(&stop), syscall.SIGTERM, SysHandler)

	// Define command: command-line arg, system signal and handler
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, SysHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, SysHandler)

	// Define daemon context
	dmn := &daemon.Context{
		PidFileName: "/var/run/" + utils.App_Tag + ".pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/" + utils.App_Tag + ".log",
		LogFilePerm: 0640,
		WorkDir:     "/",
		Umask:       027,
	}

	// Send commands if needed
	if len(daemon.ActiveFlags()) > 0 {
		d, err := dmn.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the "+utils.App_Name+" daemon:", err)
		}
		daemon.SendCommands(d)
		fmt.Println(errors.New("Signal sent!"))
		return nil
	}
	return dmn
}

func handleTaskRun() {

	if len(*task_id) > 0 {
		loadDatabase()
		_task, er := database.FindTaskById(*task_id)
		if er == nil {
			database.RunTask(_task)
		} else {
			fmt.Println("Error:", er)
		}
	}

}

func handleCompile() {

	if !*compile {
		return
	}
	if len(*config_file) > 0 {
		loadDatabase()
		database.InitConfigs()
		database.SaveTasks()
	} else {
		fmt.Println("conf file is required")
	}

}

func handleDatabaseQueries() {
	handleTaskRun()
	if *list_tasks {
		loadDatabase()
		for _, v := range database.Tasks {
			fmt.Println("------------------")
			fmt.Println("ID = ", v.Id)
			fmt.Println("Command = ", v.Command)
			fmt.Println("Next = ", time.Unix(v.When(), 0))
			fmt.Println("------------------")
		}
	}
}
