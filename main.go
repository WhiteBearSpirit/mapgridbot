package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type wpt struct {
	Lat  string `xml:"lat,attr"`
	Lon  string `xml:"lon,attr"`
	Name string `xml:"name"`
}

type gpx struct {
	XMLName xml.Name `xml:"gpx"`
	Xmlns   string   `xml:"xmlns,attr"`
	Creator string   `xml:"creator,attr"`
	Version string   `xml:"version,attr"`
	Wpts    []wpt    `xml:"wpt"`
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const aLen = 26

func genGrid(lat0, lon0, latMin, lonMax, gridStep float64) (xGrid []wpt) {

	latLimit := latMin > 0 && latMin < lat0
	lonLimit := lonMax > 0 && lonMax > lon0

	for i := 0; i < aLen; i++ {
		x, _ := Direct(lat0, lon0, gridStep*float64(i), -math.Pi)
		if latLimit && x < latMin {
			break
		}

		for j := 0; j < aLen; j++ {
			_, y := Direct(x, lon0, gridStep*float64(j), math.Pi/2.0)
			if lonLimit && y > lonMax {
				break
			}

			name := alphabet[i:i+1] + strconv.Itoa(j+1)
			xGrid = append(xGrid, wpt{
				Lat:  fmt.Sprintf("%.6f", x),
				Lon:  fmt.Sprintf("%.6f", y),
				Name: name})
		}
	}
	return
}

func getHelpMessage(langCode string) string {

	if langCode == "ru" {

		return `Введите широту и долготу в десятичном формате: гг.гггггг, гг.гггггг	
Первая координата (широта, долгота) - это самая северо-западная точка сетки, верхний левый угол. Если написать только одну координату - сетка будет 26х26 точек (максимальный размер).
Если написать вторую координату - это будет точка, ограничивающая сетку с юго-востока, нижний правый угол будет не дальше этой точки.
Если после этого написать ещё слово (одно слово, но может содержать дефисы и подчёркивания) - оно будет использовано в названии результирующего файла.
Последний параметр - шаг сетки в метрах, по умолчанию 100.
Параметры можно друг от друга отделять запятыми, пробелами, точкой с запятой.

Пример 1: 60.0000, 30.0000
Пример 2: 61.231715, 30.024984 61.223463, 30.051560 Хепосаари
Пример 3: 59.975964, 30.268458; 59.949592, 30.336530 Петроградка 200
Пример 4: 60.061522, 30.142147 Санкт-Петербург 1000`
	}

	return `Enter latitude and longitude in decimal format: xx.xxxxxx, yy.yyyyyy
First coordinates are for north-west (left top) corner.
Second coordinates are limiting grid size by south-east corner (optional).
Next parameter is the name of the grid (optional).
Last parameter is step size of the grid in meters (optional). Ex: 50, 100, 200, 500 etc.
Maximal grid size is 26 x 26 points.
Default grid step is 100 meters.

Example 1: 60.0000, 30.0000
Example 2: 61.231715, 30.024984 61.223463, 30.051560 Heposaari
Example 3: 59.975964, 30.268458; 59.949592, 30.336530 Petrogradsky_District 200
Example 4: 60.061522, 30.142147 Saint-Petersburg 1000`
}

func main() {

	log.Print("Starting telegram bot...")
	bot, err := tgbotapi.NewBotAPI("***")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	paramRegex := regexp.MustCompile(`([_\d\p{L}.-]+)`)
	coordRegex := regexp.MustCompile(`^([-+]?\d+(\.\d+)?)$`)

	for update := range updates {

		if update.Message == nil {
			// ignore any non-Message Updates
			continue
		}

		helpMessage := getHelpMessage(update.Message.From.LanguageCode)

		inputParams := paramRegex.FindAllString(update.Message.Text, -1)
		log.Print(inputParams)

		if len(inputParams) < 2 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
			bot.Send(msg)
			continue
		}

		lat0, err0 := strconv.ParseFloat(inputParams[0], 64)
		lon0, err1 := strconv.ParseFloat(inputParams[1], 64)

		if err0 != nil || err1 != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
			bot.Send(msg)
			continue
		}

		latMin, lonMax := 0.0, 0.0
		gridName := "grid"
		gridStep := 100.0
		hasSecondCoord := len(inputParams) >= 4 && coordRegex.MatchString(inputParams[2]) && coordRegex.MatchString(inputParams[3])
		if hasSecondCoord {
			latMin, _ = strconv.ParseFloat(inputParams[2], 64)
			lonMax, _ = strconv.ParseFloat(inputParams[3], 64)

			if len(inputParams) >= 5 {
				gridName = inputParams[4]
			}
			if len(inputParams) == 6 {
				gridStep, err = strconv.ParseFloat(inputParams[5], 64)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
					bot.Send(msg)
					continue
				}
			}
		} else {
			if len(inputParams) >= 3 {
				gridName = inputParams[2]
			}
			if len(inputParams) == 4 {
				gridStep, err = strconv.ParseFloat(inputParams[3], 64)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
					bot.Send(msg)
					continue
				}
			}
		}

		xGrid := genGrid(lat0, lon0, latMin, lonMax, gridStep)
		gpxData := &gpx{
			Xmlns:   "http://www.topografix.com/GPX/1/1",
			Creator: "@MapGridBot",
			Version: "1.1",
			Wpts:    xGrid}

		buf := new(bytes.Buffer)
		buf.WriteString(xml.Header)
		enc := xml.NewEncoder(buf)
		if err = enc.Encode(gpxData); err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		intStep := int(gridStep)
		fileBytes := tgbotapi.FileBytes{
			Name:  gridName + "_" + fmt.Sprint(intStep) + ".gpx",
			Bytes: buf.Bytes(),
		}

		msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, fileBytes)
		bot.Send(msg)
	}
}
