package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// RunBot ...
func RunBot(config *Config) {
	f, err := os.OpenFile("bot19rus.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("==============================================\n ~~~~~~ Bot19rus log is started ~~~~~~")
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	isUpdate := make(chan bool)
	var diff int
	go updater(isUpdate, &diff, config)
	go notificate(isUpdate, *bot, &diff)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		tgid := update.Message.From.ID

		if update.Message.Text == "test" {
			log.Println("Received test command. Processing update.")
			isUpdate <- true
			log.Println("isUpdate channel setted in true.")
		}

		//check user for exist
		if IsExist(tgid) {
			log.Println("Reseived some message. Sending last \"Cases\".")
			nums, err := GetLast()
			if err != nil {
				log.Panic(err)
			}
			numsCases := nums.Cases
			m := fmt.Sprintf("Количество заражений не менялось: %v", strconv.Itoa(numsCases))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, m)
			bot.Send(msg)
		} else {
			log.Printf("Add new user with id [%v], to database", update.Message.From.ID)
			UserAdd(tgid)
			nums, err := GetLast()
			if err != nil {
				log.Panic(err)
			}
			numsCases := nums.Cases
			m := fmt.Sprintf("Отлично. Когда общее количество заражений изменится, вы получите от меня уведомление. Сейчас общее количество заражений: %v", numsCases)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, m)
			bot.Send(msg)
		}
	}
}

func updater(isUpdate chan bool, diff *int, config *Config) {
	for {
		scrappedCases, err := ScrapeNums(config)
		log.Printf("Received number of cases from scraper: %v", scrappedCases)
		if err != nil {
			panic(err)
		}
		dbNums, err := GetLast()
		dbCases := dbNums.Cases
		log.Printf("Number of cases in database: %v", dbCases)
		if scrappedCases != dbCases {
			log.Printf("Numbers is different. Processing update.")
			NumsAdd(scrappedCases)
			*diff = scrappedCases - dbCases
			log.Printf("diff in updater: %v", *diff)
			isUpdate <- true
			log.Println("isUpdate true sended.")
		} else {
			log.Println("The numbers is same. Nothing to do.")
		}
		time.Sleep(1 * time.Hour)
		log.Println("Wait time is estimate. Continue...")
	}
}

func notificate(isUpdate chan bool, bot tgbotapi.BotAPI, diff *int) chan bool {
	for {
		if <-isUpdate == true {
			log.Println("Received update command")
			users, err := Users()
			if err != nil {
				log.Panic(err)
			}
			nums, err := GetLast()
			if err != nil {
				log.Panic(err)
			}
			numsCases := nums.Cases
			log.Printf("diff in notificator: %v", *diff)
			for _, user := range users {
				m := fmt.Sprintf("Новое число заражений: %v (+%v)", strconv.Itoa(numsCases), *diff)
				msg := tgbotapi.NewMessage(int64(user.Tgid), m)
				bot.Send(msg)
			}
			log.Println("All messages has sent.")
		}
	}
}
