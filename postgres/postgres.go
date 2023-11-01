package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rishavmngo/chat-backend-v2/entities"
	"log"
)

type Postgres struct {
	db *sql.DB
}

func InitilizePostgresStore(user, password, dbname, port string) *Postgres {
	postgres := Postgres{}
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", user, password, dbname, port)
	var err error
	postgres.db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Datbase connection established")
	// ensureTableExist(postgres.db)
	return &postgres
}

func (postgres *Postgres) GetUser(user *entities.User) error {

	err := postgres.db.QueryRow("SELECT id,created_at, last_active_at FROM users WHERE username=$1 AND password=$2", user.Username, user.Password).Scan(&user.UID, &user.CreatedAt, &user.LastActiveAt)

	if err != nil {
		return err
	}

	return nil

}

func (postgres *Postgres) GetUserByUID(user *entities.User) error {

	err := postgres.db.QueryRow("SELECT id,username,created_at, last_active_at FROM users WHERE id=$1", user.UID).Scan(&user.UID, &user.Username, &user.CreatedAt, &user.LastActiveAt)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}
func (postgres *Postgres) GetAllChatsByUID(UID uint) ([]*entities.Chats, error) {
	// var chats entities.Chats
	rows, err := postgres.db.Query(getAllChatsByUIDQuery, UID)

	if err != nil {
		log.Println("error occured", err)
		return nil, err
	}

	defer rows.Close()
	var chats []*entities.Chats
	for rows.Next() {
		var chat entities.Chats
		if err := rows.Scan(&chat.ID, &chat.CreatedAt, &chat.LastActiveAt, &chat.PartnerName, &chat.PartnerUID); err != nil {
			log.Println(err)
			return nil, err
		}
		chats = append(chats, &chat)
	}

	return chats, nil
}

func (postgres *Postgres) GetAllMsgByChatID(ChatID uint) ([]*entities.Message, error) {
	// var chats entities.Chats
	rows, err := postgres.db.Query(getAllMsgByIDQuery, ChatID)

	if err != nil {
		log.Println("error occured", err)
		return nil, err
	}

	defer rows.Close()
	var messages []*entities.Message
	for rows.Next() {
		var message entities.Message
		if err := rows.Scan(&message.ID, &message.Content, &message.ChatID, &message.CreatedAt, &message.AuthorID, &message.AuthorName); err != nil {
			log.Println(err)
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func (postgres *Postgres) GetOtherMember(chat_id, partner_id uint) (uint, error) {

	var member_id uint
	err := postgres.db.QueryRow(getOtherMemberPChat, chat_id, partner_id).Scan(&member_id)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return member_id, nil

}

func (postgres *Postgres) AddMsgToPChat(message *entities.Message) error {

	err := postgres.db.QueryRow(insertMsgToPChatQuery, message.Content, message.ChatID, message.AuthorID).Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
