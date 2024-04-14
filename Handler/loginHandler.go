package loginHandler

import (
	"database/sql"
	"fmt"
	models "go-apis/Models"
	util "go-apis/Util"
	dbConnection "go-apis/connection"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user models.RefUser
	var userId uint

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Need to check password is correct or not

	err := dbConnection.DB.QueryRow("select * from RefUser Where username=$1", user.Username).Scan(&user.RefUserId, &user.Username, &user.Password, &user.Is_Admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while connection established to DB"})
	}
	userId = user.RefUserId
	tokenString, err := util.CreateToken(userId)
	if err != nil {
		fmt.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while CreateToken"})
		return
	}
	user.Token = tokenString
	c.JSON(http.StatusOK, gin.H{"data": user})

}

func Register(c *gin.Context) {
	var user models.RefUser
	var userId uint
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	err := dbConnection.DB.QueryRow("select RefUserId from RefUser Where username=?", user.Username).Scan(&userId)
	if err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusCreated, gin.H{"message": "username is already exist!"})
			return
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "Error while execute query fro username is already exist or not"})
		}
	}
	tokenString, err := util.CreateToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while CreateToken"})
		return
	}
	user.Token = tokenString
	c.JSON(http.StatusCreated, gin.H{"data": user})
}
