package utils

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
)

func RunDaemon(args []string) error {
	signal := flag.String("s", "", "send signal to  "+App_Name+" daemon")

	handler := func(sig os.Signal) error {
		fmt.Println("signal:", sig)
		if sig == syscall.SIGTERM {
			return daemon.ErrStop
		}
		return nil
	}

	// Define command: command-line arg, system signal and handler
	daemon.AddCommand(daemon.StringFlag(signal, "term"), syscall.SIGTERM, handler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, handler)
	flag.Parse()

	// Define daemon context
	dmn := &daemon.Context{
		PidFileName: "/var/run/" + App_Tag + ".pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/" + App_Tag + ".log",
		LogFilePerm: 0640,
		WorkDir:     "/",
		Umask:       027,
		Args:        args,
	}

	// Send commands if needed
	if len(daemon.ActiveFlags()) > 0 {
		d, err := dmn.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the "+App_Name+" daemon:", err)
		}
		daemon.SendCommands(d)
		return errors.New("Daemon ActiveFlags")
	}

	// Process daemon operations - send signal if present flag or daemonize
	child, err := dmn.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if child != nil {
		return errors.New("Daemon Child Active")
	}
	defer dmn.Release()

	go func() {
		for {
			time.Sleep(10)
		}
	}()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	return nil
}
