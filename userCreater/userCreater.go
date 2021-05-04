package main

import (
	"fmt"
	"log"
	"os"

	"bufio"

	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

//ID・ハッシュ化したメールアドレス・ハッシュ化したパスワード、秘密鍵を自動生成
func main() {
	// DB, err := sql.Open("mysql", "MineDragon:@/cacyberdojo")
	// dbMap := gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	//dbHandler, err := dbMap.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	UUID, _ := uuid.NewUUID()
	id := UUID.String()

	// b, _ := hex.DecodeString(id)
	// privateKey := ed25519.PrivateKey(b)

	scanner := bufio.NewScanner(os.Stdin)
	var mailAddress string
	if scanner.Scan() {
		mailAddress = scanner.Text()
	}
	var passWord string
	if scanner.Scan() {
		passWord = scanner.Text()
	}

	hashedMailAddress, err := bcrypt.GenerateFromPassword([]byte(mailAddress), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passWord), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ID : " + id)
	fmt.Println("mailAddress : " + string(hashedMailAddress))
	fmt.Println("passWord : " + string(hashedPassword))
	//fmt.Println("Private Key : " + string(privateKey))

}
