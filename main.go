package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type User struct {
	name     string
	password string
	rights   map[string]string
}

var (
	userOne = User{
		name:   "admin",
		rights: map[string]string{"one.txt": "111", "two.txt": "111", "three.txt": "111", "four.txt": "111", "five.txt": "111"},
	}
	userTwo = User{
		name:   "Robert",
		rights: map[string]string{"one.txt": "101", "two.txt": "000", "three.txt": "011", "four.txt": "110", "five.txt": "011"},
	}
	userThree = User{
		name:   "Guest",
		rights: map[string]string{"one.txt": "011", "two.txt": "101", "three.txt": "101", "four.txt": "000", "five.txt": "100"},
	}

	users       = []User{userOne, userTwo, userThree}
	currentUser User
	currentFile string
	//quit        = ""
)

func main() {
	for {
		menu()
	}
}

func menu() {
	var user User
	var ok = false
	// authentication loop
	for !ok {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Для доступа введите логин: ")
		scanner.Scan()
		username := scanner.Text()

		if user, ok = findUser(username); !ok {
			fmt.Println("Пользователь не найден, попробуйте ещё раз!")
			return
		}
		currentUser = user
		printRights()
		ok = true

		// selecting file loop
		ok = false
		for !ok {
			if ok = chooseFile(); !ok {
				fmt.Println("Такого файла не существует, попробуйте еще раз!")
				return
			}
			fmt.Printf("Файл выбран, Вы работаете с файлом %s\n", currentFile)
			ok = true

			// working with file loop
			ok = false
			for !ok {
				if ok = workWithFile(); !ok {
					fmt.Println("У Вас недостаточно прав!")
					return
				}
			}
		}
	}
}

func workWithFile() bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Выберите цифру действия: \n" +
		"\t1. Чтение\n" +
		"\t2. Запись\n" +
		"\t3. Передача прав\n")
	scanner.Scan()
	action, _ := strconv.Atoi(scanner.Text())
	switch action {
	case 1:
		if checkRight(action) {
			read()
		} else {
			return false
		}
	case 2:
		if checkRight(action) {
			write()
		} else {
			return false
		}
	case 3:
		if checkRight(action) {
			grant()
		} else {
			return false
		}
	}
	return true
}

func grant() {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Print("Кому хотите передать права? Введите его логин: ")
	sc.Scan()
	granted := sc.Text()
	user, ok := findUser(granted)
	if ok {
		user.rights[currentFile] = currentUser.rights[currentFile]
	}
}

func write() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Что хотите записать в файл?: ")
	scanner.Scan()
	input := scanner.Text()

	file, err := os.OpenFile(currentFile, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			return
		}
	}()

	_, err = file.WriteString(input)
	if err != nil {
		return
	}
	fmt.Println("Ваш текст успешно сохранен")
}

func read() {
	fmt.Printf("Вы выбрали Чтение файла %s, ниже представлено его содержимое: \n", currentFile)
	file, err := os.Open(currentFile)
	if err != nil {
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			return
		}
	}()
	b, err := ioutil.ReadAll(file)
	fmt.Println(string(b))
}

func checkRight(action int) bool {
	for key, element := range currentUser.rights {
		if key == currentFile && action == 1 {
			if element == "100" || element == "101" || element == "110" || element == "111" {
				return true
			}
		}
		if key == currentFile && action == 2 {
			if element == "010" || element == "011" || element == "110" || element == "111" {
				return true
			}
		}
		if key == currentFile && action == 3 {
			if element == "001" || element == "011" || element == "101" || element == "111" {
				return true
			}
		}
	}
	return false
}

func chooseFile() bool {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите название файла с которым хотите работать: ")
	sc.Scan()
	fileName := sc.Text()

	file, err := os.Open(fileName)
	if err != nil {
		return false
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}(file)

	currentFile = fileName

	return true
}

func printRights() {
	fmt.Println("Перечень ваших прав:")

	var sortedNames []string
	for key, _ := range currentUser.rights {
		sortedNames = append(sortedNames, key)
	}

	for _, fileName := range sortedNames {
		fmt.Printf("\t%s: %s\n", fileName, formatRight(currentUser.rights[fileName]))
	}
}

func findUser(username string) (User, bool) {
	for _, user := range users {
		if user.name == username {
			return user, true
		}
	}
	return User{}, false
}

func formatRight(right string) string {
	switch right {
	case "000":
		return "Полный запрет"
	case "001":
		return "Передача прав"
	case "010":
		return "Запись"
	case "011":
		return "Запись, Передача прав"
	case "100":
		return "Чтение"
	case "101":
		return "Чтение, Передача прав"
	case "110":
		return "Чтение, Запись"
	case "111":
		return "Полный доступ"
	}
	return ""
}

//}
//
//func workWithFile(fileName string) {
//	var action int
//	fmt.Printf("Файл найден, Вы работаете с %s, что Вы хотите выполнить (выберите цифру)?\n\t1. Чтение\n\t2. Запись\n\t3. Передачу прав")
//	_, err := fmt.Scan(&action)
//	if err != nil {
//		return
//	}
//	switch action {
//	case 1:
//		readFile(fileName)
//	case 2:
//		writeToFile(fileName)
//	case 3:
//		grantRights(fileName)
//	}
//}
//
//func selectFile() string {
//	var fileName string
//	fmt.Printf("Выберите файл с которым хотите работать, введите его полное название: ")
//	_, err := fmt.Scan(&fileName)
//	if err != nil {
//		return ""
//	}
//	files, err := ioutil.ReadDir("/files/")
//	if err != nil {
//		log.Fatal(err)
//	}
//	var names []string
//	for _, f := range files {
//		names = append(names, f.Name())
//	}
//	for _, name := range names {
//		if name == fileName {
//			return fileName
//		}
//	}
//	return ""
//}
//
//func grantRights(fileName string) {
//	var username string
//	fmt.Printf("Вы выбрали передачу прав. Кому хотите передать свои права? Введите имя пользователя: ")
//	_, err := fmt.Scan(&username)
//	if err != nil {
//		return
//	}
//}
//
//func writeToFile(fileName string) {
//
//}
//
//func readFile(fileName string) {
//
//}
//
//func printRightsMatrix(user *User) {
//	fmt.Printf("Вам выдан следующий перечень прав:\n")
//	keys := make([]string, 0)
//	for k := range user.rights {
//		keys = append(keys, k)
//	}
//	sort.Strings(keys)
//	for _, k := range keys {
//		fmt.Printf("\t%s: %s\n", k, formatRight(user.rights[k]))
//	}
//}
//
//func formatRight(right string) string {
//	switch right {
//	case "000":
//		return "Полный запрет"
//	case "001":
//		return "Передача прав"
//	case "010":
//		return "Запись"
//	case "011":
//		return "Запись, Передача прав"
//	case "100":
//		return "Чтение"
//	case "101":
//		return "Чтение, Передача прав"
//	case "110":
//		return "Чтение, Запись"
//	case "111":
//		return "Полный доступ"
//	default:
//		return "Ошибка прав доступа"
//	}
//}
//
//func greeting() (string, string) {
//	var name, password string
//	fmt.Print("Введите логин: ")
//	_, _ = fmt.Scan(&name)
//	fmt.Print("Введите пароль: ")
//	_, _ = fmt.Scan(&password)
//	return name, password
//}
//
//func auth(name, password string) *User {
//	for _, user := range users {
//		if user.name == name && user.password == password {
//			fmt.Printf("Вход выполнен, добро пожаловать, %s.\n", user.name)
//			return &user
//		}
//	}
//	return nil
//}
