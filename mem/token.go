package mem

import ()

func SaveToken(token string, id int) {
	token_map.Set(token, id)
}

func GetToken(token string) (int, bool) {
	var id int
	ret, ok := token_map.Get(token)

	if ok {
		id = ret.(int)
		return id, true
	} else {
		return 0, false
	}

}
