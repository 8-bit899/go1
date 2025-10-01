package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/io893/calendar_app/calendar"
	"github.com/io893/calendar_app/events"
	"github.com/io893/calendar_app/logger"
	"github.com/io893/calendar_app/storage"
)

var mu sync.Mutex

type Cmd struct {
	calendar *calendar.Calendar
	log      *Log
	logger   *logger.Logger
}
type Log struct {
	msg     []string
	storage storage.Store
}

func NewCmd(c *calendar.Calendar, s storage.Store, logName string) (*Cmd, error) {
	logger, err := logger.LoggerNew(logName)
	if err != nil {
		return nil, fmt.Errorf("ошибка cmd: %w", err)
	}
	return &Cmd{
		calendar: c,
		log:      NewLog(s),
		logger:   logger,
	}, nil
}
func NewLog(s storage.Store) *Log {
	return &Log{
		msg:     []string{},
		storage: s,
	}
}
func (l *Log) Logsave() error {

	data, err := json.Marshal(l.msg)
	if err != nil {

		return err
	}
	err = l.storage.Save(data)
	return err
}
func (l *Log) Logload() error {

	data, err := l.storage.Load()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &l.msg)
	if err != nil {

		return err
	}
	return err
}
func (c *Cmd) Logread() {

	for _, msg := range c.log.msg {
		c.calendar.Notify(msg)
	}

}
func (l *Log) Logwrite(msg string) {
	mu.Lock()
	defer mu.Unlock()
	l.msg = append(l.msg, msg)
}
func (c *Cmd) executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	c.log.Logwrite("cmd: " + input)
	c.logger.Info(input)
	parts, err := shlex.Split(input)

	if err != nil {
		fmt.Println(err)
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "logger":
	case "log":
		c.Logread()
	case "add":
		if len(parts) < 4 {
			err := "Формат: add \"название события\" \"дата и время\" \"приоритет\""
			c.calendar.Notify(err)
			c.logger.Error(err)
			return
		}
		title := parts[1]
		date := parts[2]
		priority := events.Priority(parts[3])
		e, err := c.calendar.AddEvent(title, date, priority)
		if err != nil {
			descriptErr := "Ошибка добавления: " + err.Error()
			c.calendar.Notify(descriptErr)
			c.logger.Error(descriptErr)
		} else {
			descriptAction := "Событие: " + e.Title + " добавлено"
			c.calendar.Notify(descriptAction)
			c.logger.Info(descriptAction)
		}
	case "update":
		if len(parts) < 5 {
			c.calendar.Notify("Формат: update \"id события\" \"название события\" \"дата и время\" \"приоритет\"")
			return
		}
		key := parts[1]
		title := parts[2]
		date := parts[3]
		priority := events.Priority(parts[4])
		c.calendar.EditEvent(key, title, date, priority)

	case "remove":
		if len(parts) < 2 {
			c.calendar.Notify("Формат: remove \"id события\"")
			return
		}
		key := parts[1]

		c.calendar.DeleteEvent(key)

	case "list":
		c.calendar.Notify("Список событий календаря:")
		c.calendar.ShowEvents()
		fmt.Println()
	case "reminder":
		if len(parts) < 4 {
			c.logger.Error("Формат: Reminder \"id события\" \"сообщение\" \"дата и время\" ")
			c.calendar.Notify("Формат: Reminder \"id события\" \"сообщение\" \"дата и время\" ")

			return
		}
		key := parts[1]
		msg := parts[2]
		date := parts[3]
		err := c.calendar.SetEventReminder(key, msg, date)
		if err != nil {
			c.calendar.Notify(err.Error())
			c.logger.Error(err.Error())
		}

	case "reminderRmv":
		if len(parts) < 2 {
			c.calendar.Notify("Формат: Reminder \"id события\"")
			return
		}
		key := parts[1]
		err := c.calendar.CancelEventReminder(key)
		if err != nil {
			c.logger.Error(err.Error())
			c.calendar.Notify(err.Error())
		}
	case "help":

		c.calendar.Notify("add - Добавить событие.\n Для добавления события введите команду в формате: add \"название события\" \"дата и время\" \"приоритет\" ")

		c.calendar.Notify("list - Выводит список всех событий календаря. Команда используется без параметров. ")

		c.calendar.Notify("remove - Удалить событие.\n Для удаления события введите команду в формате: remove \"id события\" ")

		c.calendar.Notify("update - Изменить событие.\n Для изменения события введите команду в формате:  update \"id события\" \"название события\" \"дата и время\" \"приоритет\"")

		c.calendar.Notify("reminder - Устанавливает уведомление для события.\n Для добавления события введите команду в формате: reminder \"id события\" \"сообщение\" \"дата и время\" ")

		c.calendar.Notify("reminderRmv - Удаляет уведомление для события.\n Для добавления события введите команду в формате: reminderRmv \"id события\" ")

		c.calendar.Notify("log - выводит действия пользователя")

		c.calendar.Notify("exit - Выйти из программы. Команда используется без параметров")

		c.calendar.Notify("help - Показать справку")

	case "exit":
		err := c.log.Logsave()
		if err != nil {
			fmt.Println("Ошибка сохранения лога:", err)
		}
		c.calendar.Save()
		c.calendar.CloseNotify()

	default:
		c.calendar.Notify("Неизвестная команда:")
		c.calendar.Notify("Введите 'help' для просмотра списка команд")
	}
}
func (c *Cmd) exitChecker(in string, breakline bool) bool {
	if !breakline {
		return false
	}
	input := strings.TrimSpace(in)
	return input == "exit"
}
func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие"},
		{Text: "list", Description: "Показать все события"},
		{Text: "remove", Description: "Удалить событие"},
		{Text: "help", Description: "Показать справку"},
		{Text: "log", Description: "Вывести лог программы"},
		{Text: "exit", Description: "Выйти из программы"},
		{Text: "reminder", Description: "Установить уведомление"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
func (c *Cmd) Run() {
	defer func() {
		cmd := exec.Command("stty", "sane")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}()
	err := c.log.Logload()
	if err != nil {
		fmt.Println("Файла с записями лога нет. Создали новый:", err)
	}
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionSetExitCheckerOnInput(c.exitChecker),
	)

	go func() {
		for msg := range c.calendar.Notification {
			fmt.Println(msg)
			c.logger.Info(msg)
			c.log.Logwrite(msg)
		}
	}()
	p.Run()
}
