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
	err := dbConnection.DB.QueryRow("select RefUserId from RefUser Where username=$1", user.Username).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			goto createToken
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while execute query from username is already exist or not"})
			return
		}
	}
	if userId != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username is already exist"})
		return
	}

createToken:
	tokenString, err := util.CreateToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while CreateToken"})
		return
	}
	user.Token = tokenString
	//Need to add encypt password
	fmt.Println(user.Password)

	insertStmt := `insert into RefUser ("username", "userpassword","isadmin") values($1, $2, $3)`
	res, e := dbConnection.DB.Exec(insertStmt, user.Username, user.Password, user.Is_Admin)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
		return
	}
	_, fail := res.RowsAffected()
	if fail != nil {
		fmt.Println(fail)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while Insert query for new user"})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"data": user})
		return
	}
}
