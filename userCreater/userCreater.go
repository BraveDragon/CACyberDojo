package main

import (
	"CACyberDojo/model/usermodel"
	"fmt"
	"log"
	"os"

	"bufio"

	"database/sql"

	"github.com/go-gorp/gorp"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

//ID・ハッシュ化したメールアドレス・ハッシュ化したパスワード、秘密鍵を自動生成してDBに格納
func main() {
	DB, err := sql.Open("mysql", "MineDragon:@/cacyberdojo")
	if err != nil {
		fmt.Println("Error occurred when connecting to MySQL.")
		log.Fatal(err)
	}
	dbMap := &gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	//DBのテーブルと構造体を結びつける
	dbMap.AddTableWithName(usermodel.User{}, "users")
	dbHandler, err := dbMap.Begin()
	if err != nil {
		fmt.Println("Error occurred when creating DBMap.")
		log.Fatal(err)
	}

	UUID, _ := uuid.NewUUID()
	id := UUID.String()

	// b, _ := hex.DecodeString(id)
	// privateKey := ed25519.PrivateKey(b)

	scanner := bufio.NewScanner(os.Stdin)
	var name string
	if scanner.Scan() {
		name = scanner.Text()
	}
	var mailAddress string
	if scanner.Scan() {
		mailAddress = scanner.Text()
	}
	var passWord string
	if scanner.Scan() {
		passWord = scanner.Text()
	}

	// hashedMailAddress, err := bcrypt.GenerateFromPassword([]byte(mailAddress), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passWord), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("ID : " + id)
	fmt.Println("name : " + name)
	fmt.Println("mailAddress : " + mailAddress)
	fmt.Println("passWord : " + passWord)

	err = dbHandler.Insert(&usermodel.User{
		Id:          id,
		Name:        name,
		PassWord:    passWord,
		MailAddress: mailAddress,
		//PrivateKey:  privateKey,
	})
	if err != nil {
		fmt.Println("Error occurred when inserting.")
		log.Fatal(err)
	}
	err = dbHandler.Commit()
	if err != nil {
		fmt.Println("Error occurred when committing.")
		log.Fatal(err)
	}

}
