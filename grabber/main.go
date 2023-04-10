package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
)

type ProgramParameters struct {
	source            string
	destinationFolder string
}

func main() {
	parameters := getProgramParameters()
	isValid := validateParameters(parameters)
	if !isValid {
		log.Fatal("Paramers validation failed")
	}
	urls, err := readSitesFromFile(parameters.source)
	if err != nil {
		log.Fatal("Failed to read urls: ", err)
	}

	parallelLoadPagesToDestinationDirecory(urls, parameters.destinationFolder)
	fmt.Println("Done")
}

// Получение параметров программы
// source - Путь к файлу со списком url сайтов для загрузки
// destination - Путь к дерриктории в которую будут сохранены сайты
func getProgramParameters() ProgramParameters {
	source := flag.String("source", "source.txt", "Путь к файлу со списком url сайтов для загрузки")
	destinationFolder := flag.String("destination", "dest", "Путь к дерриктории в которую будут сохранены сайты")

	flag.Parse()

	return ProgramParameters{source: *source, destinationFolder: *destinationFolder}
}

// Проверка валидности параметров программы, что они не nil и файл с папкой существуют
func validateParameters(parameters ProgramParameters) bool {
	sourceStat, err := os.Stat(parameters.source) // Получаю FileInfo для destinationFolder

	if err != nil || sourceStat.IsDir() { // Проверяю что не было ошибки при полчении source FileInfo и что sourceStat не папка
		log.Println("Source file not found")
		return false
	}

	folderStat, err := os.Stat(parameters.destinationFolder) // Получаю FileInfo для destinationFolder

	if err != nil || !folderStat.IsDir() { // Проверяю что не было ошибки при полчении FileInfo и что folderStat папка
		log.Println("Target folder not found")
		return false
	}

	return true
}

// Читает список сайтов из файла
func readSitesFromFile(pathToSourceFile string) ([]string, error) {
	file, err := os.Open(pathToSourceFile)

	if err != nil {
		log.Printf("Error opening file " + pathToSourceFile)
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	return urls, nil
}

// Загружает все страницы по списку urls и записывает их в папку destination
func parallelLoadPagesToDestinationDirecory(urls []string, destination string) {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		url := url
		go func(url string) { // создаю анонимную функцию и вызываю ее в горутине
			defer wg.Done()
			body, err := loadPage(url)
			if err == nil {
				writeToFile(fmt.Sprintf("%s/%s", destination, createFileNameFromUrl(url)), body)
			}
		}(url)
	}
	wg.Wait()
}

// Преобразует url в имя файла
func createFileNameFromUrl(url string) string {
	httpRegex := regexp.MustCompile(`http(s)?:\/\/`)
	pointRegex := regexp.MustCompile(`[.\/]`)
	url = httpRegex.ReplaceAllString(url, "")
	url = pointRegex.ReplaceAllString(url, "_")
	return url + ".html"
}

// Загружает страницу по url и возвращает тело
func loadPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error loading page " + url)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body " + url)
		return nil, err
	}
	return body, nil
}

// Создает и записывает в файл content
func writeToFile(path string, content []byte) error {
	file, err := os.Create(path)

	if err != nil {
		log.Println("Failed to create file " + path)
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		log.Println("Failed to write to file", err)
		return err
	}
	log.Println("Loaded " + path)
	return nil
}
