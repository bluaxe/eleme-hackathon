package cache

type keys struct {
	TOKEN_KEY string
}

var k keys = InitKeys()

func InitKeys() keys {
	var k keys
	k.TOKEN_KEY = `tok:%s`
	return k
}
