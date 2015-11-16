package mem

import ()

func SaveToken(token string, id int) {
	token_map[token] = id
}

func GetToken(token string) (int, bool) {
	id, ok := token_map[token]
	return id, ok
}

func UserGetToken(uid int) string {
	return uid_token_map[uid]
}

func UserSetToken(uid int, token string) {
	uid_token_map[uid] = token
}
