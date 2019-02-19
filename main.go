package main

var defaultPaths = []string{
	"/etc/gitconfig",
	"~/.gitconfig",
	"~/.config/git/config",
	".git/config",
}

func main() {
	// todo: Add support of gui
	InitCli()
//	filePaths := make(chan string)
//	userEmail := make(chan []byte)
//	userNames := make(chan []byte)
//	finish := make(chan struct{})
//	reader := bufio.NewReader(os.Stdin)
//
//	go GetUserNamesAndEmail(filePaths, userEmail, userNames, finish)
//
//	for _, defaultPath := range defaultPaths {
//		filePaths <- defaultPath
//	}
//	close(filePaths)
//	uNames := make([][]byte, 0, 0)
//	uEmail := make([][]byte, 0, 0)
//LOOP:
//	for {
//		select {
//		case u := <-userNames:
//			uNames = append(uNames, u)
//		case e := <-userEmail:
//			uEmail = append(uEmail, e)
//		case <-finish:
//			break LOOP
//		default:
//		}
//	}
//	fmt.Println("===Users Names====")
//	for _, userName := range uNames {
//		fmt.Printf("%q\n", userName)
//	}
//
//	fmt.Println("===Users Emails====")
//	for _, email := range uEmail {
//		fmt.Printf("%q\n", email)
//	}
//
//	fmt.Print("Type User Name from the list:")
//	chUserName, _ := reader.ReadString('\n')
//
//	fmt.Print("Type User Email from the list:")
//	chUserEmail, _ := reader.ReadString('\n')
//
//	user := User{
//		[]byte(chUserName),
//		[]byte(chUserEmail),
//	}
//
//	err := UpdateUserInfo("test_file", []byte(user.UserRepresentation()))
//	fmt.Println(err)
}
