package main

import (
	"fmt"

	"github.com/io893/calendar_app/calendar"
	"github.com/io893/calendar_app/cmd"
	"github.com/io893/calendar_app/storage"
)

func main() {
	//	s := storage.NewJsonStorage("calendar.json")

	zs := storage.NewZipStorage("calendar.zip")
	logjs := storage.NewJsonStorage("log.json")

	c := calendar.NewCalendar(zs)
	err := c.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки данных:", err)
		return
	}

	cli := cmd.NewCmd(c, logjs)
	cli.Run()
	fmt.Println("=========================================================================")
}
