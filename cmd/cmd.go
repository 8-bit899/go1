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
	c.HandleCommand(cmd, parts)
}

func (c *Cmd) notifyUserAndLogError(message string) {
	c.calendar.Notify(message)
	c.logger.Error(message)
}

func (c *Cmd) notifyUserAndLogInfo(message string) {
	c.calendar.Notify(message)
	c.logger.Info(message)
}

func (c *Cmd) validateArgs(parts []string, requiredCount int, format string) bool {
	if len(parts) < requiredCount {
		errMsg := fmt.Sprintf("Неверное количество аргументов. Формат: %s", format)
		c.notifyUserAndLogError(errMsg)
		return false
	}
	return true
}

func (c *Cmd) handleEventAdd(operation string, parts []string) {

	if !c.validateArgs(parts, 4, fmt.Sprintf("%s \"название\" \"дата\" \"приоритет\"", operation)) {
		return
	}

	title := parts[1]
	date := parts[2]
	priority := events.Priority(parts[3])

	e, err := c.calendar.AddEvent(title, date, priority)
	if err != nil {
		c.notifyUserAndLogError(fmt.Sprintf("Ошибка %s: %v", operation, err))
	} else {
		c.notifyUserAndLogInfo(fmt.Sprintf("Событие '%s' %s", e.Title, getOperationPastTense(operation)))
	}
}

func (c *Cmd) handleEventUpdate(parts []string) {
	if !c.validateArgs(parts, 5, "update \"id\" \"название\" \"дата\" \"приоритет\"") {
		return
	}

	key := parts[1]
	title := parts[2]
	date := parts[3]
	priority := events.Priority(parts[4])

	err := c.calendar.EditEvent(key, title, date, priority)
	if err != nil {
		c.notifyUserAndLogError(fmt.Sprintf("Ошибка обновления: %v", err))
	} else {
		c.notifyUserAndLogInfo(fmt.Sprintf("Событие с ID %s обновлено", key))
	}
}

func (c *Cmd) handleEventRemoval(parts []string) {
	if !c.validateArgs(parts, 2, "remove \"id события\"") {
		return
	}

	key := parts[1]
	err := c.calendar.DeleteEvent(key)
	if err != nil {
		c.notifyUserAndLogError(fmt.Sprintf("Ошибка удаления: %v", err))
	} else {
		c.notifyUserAndLogInfo(fmt.Sprintf("Событие с ID %s удалено", key))
	}
}

func (c *Cmd) handleReminderSet(parts []string) {
	if !c.validateArgs(parts, 4, "reminder \"id\" \"сообщение\" \"дата\"") {
		return
	}

	key := parts[1]
	msg := parts[2]
	date := parts[3]

	err := c.calendar.SetEventReminder(key, msg, date)
	if err != nil {
		c.notifyUserAndLogError(fmt.Sprintf("Ошибка установки напоминания: %v", err))
	} else {
		c.notifyUserAndLogInfo(fmt.Sprintf("Напоминание для события %s установлено", key))
	}
}

func (c *Cmd) handleReminderRemoval(parts []string) {
	if !c.validateArgs(parts, 2, "reminderRmv \"id события\"") {
		return
	}

	key := parts[1]
	err := c.calendar.CancelEventReminder(key)
	if err != nil {
		c.notifyUserAndLogError(fmt.Sprintf("Ошибка удаления напоминания: %v", err))
	} else {
		c.notifyUserAndLogInfo(fmt.Sprintf("Напоминание для события %s удалено", key))
	}
}

func (c *Cmd) showHelp() {
	helpCommands := map[string]string{
		"add":         "Добавить событие. Формат: add \"название\" \"дата\" \"приоритет\"",
		"list":        "Выводит список всех событий календаря. Без параметров.",
		"remove":      "Удалить событие. Формат: remove \"id события\"",
		"update":      "Изменить событие. Формат: update \"id\" \"название\" \"дата\" \"приоритет\"",
		"reminder":    "Установить уведомление. Формат: reminder \"id\" \"сообщение\" \"дата\"",
		"reminderRmv": "Удалить уведомление. Формат: reminderRmv \"id события\"",
		"log":         "Выводит действия пользователя",
		"exit":        "Выйти из программы. Без параметров",
		"help":        "Показать справку",
	}

	for cmd, desc := range helpCommands {
		c.calendar.Notify(fmt.Sprintf("%s - %s", cmd, desc))
	}
	c.logger.Info("Пользователь запросил справку")
}

func (c *Cmd) handleExit() {
	err := c.log.Logsave()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Ошибка сохранения лога: %v", err))
		fmt.Println("Ошибка сохранения лога:", err)
	}
	err = c.calendar.Save()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Ошибка сохранения календаря: %v", err))
	}

	c.calendar.CloseNotify()
	c.logger.Info("Приложение завершило работу")
}

func getOperationPastTense(operation string) string {
	tenses := map[string]string{
		"add":    "добавлено",
		"update": "обновлено",
		"remove": "удалено",
	}
	if tense, ok := tenses[operation]; ok {
		return tense
	}
	return "обработано"
}

func (c *Cmd) HandleCommand(cmd string, parts []string) {
	c.logger.Info(fmt.Sprintf("Выполнение команды: %s", cmd))

	switch cmd {
	case "logger":

		c.logger.Info("Выполнена команда logger")

	case "log":
		c.logger.Info("Чтение лога пользователем")
		c.Logread()

	case "add":
		c.handleEventAdd("add", parts)

	case "update":
		c.handleEventUpdate(parts)

	case "remove":
		c.handleEventRemoval(parts)

	case "list":
		c.logger.Info("Запрос списка событий")
		c.calendar.Notify("Список событий календаря:")
		c.calendar.ShowEvents()
		fmt.Println()

	case "reminder":
		c.handleReminderSet(parts)

	case "reminderRmv":
		c.handleReminderRemoval(parts)

	case "help":
		c.showHelp()

	case "exit":
		c.handleExit()

	default:
		errMsg := fmt.Sprintf("Неизвестная команда: %s", cmd)
		c.notifyUserAndLogError(errMsg)
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
		{Text: "update", Description: "Изменить событие"},
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
			c.log.Logwrite(msg)
		}
	}()
	p.Run()
}
