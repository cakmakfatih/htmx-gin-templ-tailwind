package models

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

type UserModel struct {
	ID       string
	Email    string
	MetaData UserMetaData
}

type UserMetaData struct {
	Picture  string `json:"picture"`
	FullName string `json:"full_name"`
}

func (UserModel) UserFromContext(c *gin.Context) *UserModel {
	user, ok := c.Get("user")

	if !ok {
		return nil
	}

	supabaseUser := user.(*supabase.User)
	userMetadata, err := json.Marshal(supabaseUser.UserMetadata)

	if err != nil {
		return nil
	}

	metaData := &UserMetaData{}
	json.Unmarshal(userMetadata, &metaData)

	return &UserModel{
		ID:       supabaseUser.ID,
		Email:    supabaseUser.Email,
		MetaData: *metaData,
	}
}
