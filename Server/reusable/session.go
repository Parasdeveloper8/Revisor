package reusable

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// type to store Key-value
type SessionKeyValue struct {
	Key   string
	Value string
}

type SessionFuncs interface {
	//set keys and their values in session
	SessionSet(c *gin.Context, keysValues []SessionKeyValue) error
}

func SessionSet(c *gin.Context, keysValues []SessionKeyValue) error {
	session := sessions.Default(c)
	//loop keysValues to set each pair of key-value
	for _, keyval := range keysValues {
		//set key and value
		session.Set(keyval.Key, keyval.Value)
	}

	err := session.Save()
	if err != nil {
		return err
	}
	return nil
}
