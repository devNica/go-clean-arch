package controller

import (
	"net/http"
	"strconv"

	"github.com/devNica/go-clean-arch/dto"
	"github.com/devNica/go-clean-arch/entity"
	"github.com/devNica/go-clean-arch/helper"
	"github.com/devNica/go-clean-arch/service"
	"github.com/gin-gonic/gin"
)

// AuthController interface is a contract what this controller can do
type AuthController interface {
	Signup(ctx *gin.Context)
	Signin(ctx *gin.Context)
}

type authController struct {
	// this is where you put your service
	authService service.AuthService
	jwtService  service.JWTService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Signin(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := helper.CreateErrorResponse("Failed to process requets", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken
		response := helper.CreateResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := helper.CreateErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Signup(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helper.CreateErrorResponse("Failed to process requets", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if !c.authService.IsDuplicatedEmail(registerDTO.Email) {
		response := helper.CreateErrorResponse("Failed to process requets", "Duplicated Email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.CreateResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusOK, response)
	}

}
