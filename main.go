package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	URL_BASE        = "https://api.trello.com/1/"
	ID_BOARD        = "55d5487fb929f379cf626053"
	DEBUG           = false
	NAME_FILE_KEY   = "./key"
	NAME_FILE_TOKEN = "./authorize"
)

var KEY string

func main() {
	if !fileExists(NAME_FILE_KEY) {
		log.Fatalln("File " + NAME_FILE_KEY + " not found.")
	}

	// get key
	KEY = readFile(NAME_FILE_KEY)

	// argsWithProg := os.Args
	// argsWithoutProg := os.Args[1:]

	// arg := os.Args[3]
	token := login()

	getBoard(ID_BOARD, token)
	// fmt.Println(argsWithoutProg)
	// fmt.Println(arg)
}

// connection with api trello
// method to request get

type Board struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func getBoard(idBoard string, token string) Board {
	params := "?fields=id,name,url&key=" + KEY + "&token=" + token
	prefix := "boards/"

	url := URL_BASE + prefix + idBoard + params

	result := get(url)

	board := Board{}
	json.Unmarshal([]byte(result), &board)

	fmt.Println(board)

	return board
}

func login() string {
	if !fileExists(NAME_FILE_TOKEN) {
		return wirteToken()
	}

	token := readFile(NAME_FILE_TOKEN)
	fmt.Println("User logged with token " + token)

	return token
	// b615a133ea49b3b785322b697d538ce73783d386b7b6b78612d976e8ec4ce1a4
}

func wirteToken() string {
	url := "https://trello.com/1/authorize/?expiration=never&name=GOTrello&scope=read,write&response_type=token&key=" + KEY
	// read token
	fmt.Println("Please access this url: " + url)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter with token: ")
	token, _ := reader.ReadString('\n')

	// write file
	writer, err := os.Create("./authorize")
	if err != nil {
		log.Fatalln(err)
	}

	defer writer.Close()

	w, err := writer.WriteString(token)
	if err != nil {
		log.Fatalln(err)
	}

	writer.Sync()

	if DEBUG {
		fmt.Println("Write byts " + string(w))
	}

	return strings.TrimSpace(token)
}

func fileExists(nameFile string) bool {
	info, err := os.Stat(nameFile)

	if err != nil && (os.IsNotExist(err) || info.IsDir()) {
		return false
	}

	return true
}

func readFile(nameFile string) string {
	readToken, err := ioutil.ReadFile(nameFile)
	if err != nil {
		log.Fatalln(err)
	}

	return strings.TrimSpace(string(readToken))
}

func get(url string) []byte {
	fmt.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	// to not timeout
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return body
}
