package main

import (
	"Go-Web-Scrapper-Go-DownloadModule/db"
	"Go-Web-Scrapper-Go-DownloadModule/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

var semaphore *utils.Semaphore
var wg sync.WaitGroup

func main() {
	var inputStr string

	fmt.Print("Enter Download Count At Time: ")
	_, err := fmt.Scanln(&inputStr)
	if err != nil {
		fmt.Println("[E] Error:", err)
		return
	}

	num, err := strconv.Atoi(inputStr)
	if err != nil {
		fmt.Println("[E] Error: Please enter a valid integer")
		return
	}
	semaphore = utils.NewSemaphore(num)
	db.GetFileFromDb()

	for {
		data, ok := <-db.FileChannel
		if !ok {
			fmt.Println("Data Channel closed. Programme Will Exit. 0")
			return
		}
		if strings.Contains(data.Url, "thumb_") {
			continue
		}
		wg.Add(1)
		go func() {
			semaphore.Acquire()
			downloader(&utils.File{Url: data.Url})
			data.DoneAndSave()
			semaphore.Release()
			wg.Done()

		}()
	}

}

func downloader(file *utils.File) {
	file.CreatePath()
	err := file.StartDownload()
	db.DownloadsC++
	if err != nil {
		log.Println(err)
		return
	}

}
