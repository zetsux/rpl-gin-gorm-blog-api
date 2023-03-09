package controller

import (
	"go-blogrpl/dto"
	"go-blogrpl/entity"
	"go-blogrpl/service"
	"go-blogrpl/utils"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
}

type UserController interface {
	SignUp(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetUserByUsername(ctx *gin.Context)
}

func NewUserController(userS service.UserService) UserController {
	return &userController{userService: userS}
}

func (userC *userController) SignUp(ctx *gin.Context) {
	var userDTO dto.UserSignUpRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		resp := utils.CreateResponse("Failed to process user sign up request", http.StatusBadRequest, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	newUser, err := userC.userService.CreateNewUser(ctx, userDTO)
	if err != nil {
		resp := utils.CreateResponse(err.Error(), http.StatusBadRequest, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.CreateResponse("user signed up successfully", http.StatusCreated, newUser)
	ctx.JSON(http.StatusCreated, resp)
}

func (userC *userController) GetAllUsers(ctx *gin.Context) {
	users, err := userC.userService.GetAllUsers(ctx)
	if err != nil {
		resp := utils.CreateResponse("Failed to fetch all users", http.StatusBadRequest, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp utils.Response
	if len(users) == 0 {
		resp = utils.CreateResponse("no user found", http.StatusOK, users)
	} else {
		resp = utils.CreateResponse("successfully fetched all users", http.StatusOK, users)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := userC.userService.GetUserByUsername(ctx, username)
	if err != nil {
		resp := utils.CreateResponse(err.Error(), http.StatusBadRequest, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp utils.Response
	if reflect.DeepEqual(user, entity.User{}) {
		resp = utils.CreateResponse("user not found", http.StatusOK, nil)
	} else {
		resp = utils.CreateResponse("successfully fetched user", http.StatusOK, user)
	}
	ctx.JSON(http.StatusOK, resp)
}
