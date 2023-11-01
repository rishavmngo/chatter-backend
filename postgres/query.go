package postgres

var (
	getAllChatsByUIDQuery = "select pc1.chat_id as chat_id, chats.created_at , chats.last_active_at, users.username as partner_name , member_id as partner_uid from private_chat as pc1 join (select chat_id, member_id as me from private_chat where member_id = $1) as ot on pc1.chat_id = ot.chat_id join users on users.id = pc1.member_id join chats on chats.id = pc1.chat_id where pc1.member_id != $1"
	getAllMsgByIDQuery    = "select messages.id, content, chat_id, messages.created_at as message_created_at, author_id, users.username as author_name   from messages join users on messages.author_id = users.id where chat_id = $1 order by messages.created_at;"
	getOtherMemberPChat   = "select member_id  from private_chat where chat_id=$1 and member_id != $2;"
	insertMsgToPChatQuery = "INSERT into messages(content,chat_id,author_id,created_at) values($1,$2,$3,now()) returning id, created_at"
)
