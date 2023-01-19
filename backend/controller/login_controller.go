package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wba-bc-project-05/backend/model"
	"wba-bc-project-05/library/jwt"

	"github.com/gin-gonic/gin"
)

func (p *Controller) RespError(c *gin.Context, body interface{}, status int, err ...interface{}) {
	bytes, _ := json.Marshal(body)

	fmt.Println("Request error", "path", c.FullPath(), "body", bytes, "status", status, "error", err)

	c.JSON(status, gin.H{
		"Error":  "Request Error",
		"path":   c.FullPath(),
		"body":   bytes,
		"status": status,
		"error":  err,
	})
	c.Abort()
}

type SigninParam struct {
	ID string `json:"id" bson:"id"`
	Pw string `json:"pw" bson:"pw"`
}

// Signin - 로그인 메서드
func (p *Controller) SignIn(c *gin.Context) {
	id := c.PostForm("id")
	pw := c.PostForm("pw")

	User := model.User{}
	// req := SigninParam{ID: id, Pw: pw}
	if err := p.md.SigninModel(id, pw); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	c.JSON(500, gin.H{
		"status":       500,
		"message":      "일치하는 회원이 없습니다.",
		"refreshToken": "null",
		"accessToken":  "null",
	})
	c.Next()

	refreshToken, err := jwt.CreateRefreshToken(User.Wallet)
	if err != nil {
		c.JSON(500, gin.H{
			"status":       500,
			"message":      "refreshtoken 생성 중 에러",
			"refreshToken": "null",
			"accessToken":  "null",
		})
		return
	}

	accessToken, err := jwt.CreateAccessToken(User.Wallet, User.IsManager)
	if err != nil {
		c.JSON(500, gin.H{
			"status":       500,
			"message":      "accesstoken 생성 중 에러",
			"refreshToken": refreshToken,
			"accessToken":  "null",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":       200,
		"message":      "토큰 발급 완료",
		"refreshToken": refreshToken,
		"accessToken":  accessToken,
	})
	return
}

// SignUp - 회원가입
func (p *Controller) SignUp(c *gin.Context) {
	id := c.PostForm("id")
	pw := c.PostForm("pw")
	wallet := c.PostForm("wallet")
	privatekey := c.PostForm("privatekey")
	isManager := c.PostForm("is_manager")

	req := model.User{UserID: id, Pw: pw, Wallet: wallet, PrivateKey: privatekey, IsManager: isManager}
	if err := p.md.SignUpModel(req); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	c.JSON(500, gin.H{
		"status":  500,
		"message": "회원가입 실패",
	})
	c.Next()
}

// Logout - 로그아웃
func (p *Controller) LogOut(c *gin.Context) {

	//클라에서 토큰지워주는 방법 사용

	c.JSON(200, gin.H{
		"status":       200,
		"message":      "로그아웃",
		"accessToken":  "null",
		"refreshToken": "null",
	})
	return
}
