package controller

import (
	"WBA-BC-Project-05/library/jwt"
	"net/http"
	"wba-bc-project-05/backend/model"

	"github.com/gin-gonic/gin"
)

// Signin - 로그인 메서드
func (p *Controller) SignIn(c *gin.Context) {
	id := c.PostForm("id")
	pw := c.PostForm("pw")

	req := model.User{UserID: id, Pw: pw}
	if err := p.md.SigninModel(req); err != nil {
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

	refreshToken, err := jwt.CreateRefreshToken(req.Wallet)
	if err != nil {
		c.JSON(500, gin.H{
			"status":       500,
			"message":      "refreshtoken 생성 중 에러",
			"refreshToken": "null",
			"accessToken":  "null",
		})
		return
	}

	accessToken, err := jwt.CreateAccessToken(req.Wallet, req.IsManager)
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
		"status": 200,
		"message": "토큰 발급 완료",
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
	isManager := c.PostForm("is_manager")

	req := model.User{UserID: id, Pw: pw, Wallet: wallet, IsManager: isManager}
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