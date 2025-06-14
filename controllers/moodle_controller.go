package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rizkycahyono97/moodle-api/model/web"
	"github.com/rizkycahyono97/moodle-api/services"
	"github.com/rizkycahyono97/moodle-api/utils/validation"
	"log"
	"net/http"
)

type MoodleController struct {
	moodleService services.MoodleService
}

func NewMoodleController(moodleService services.MoodleService) *MoodleController {
	return &MoodleController{moodleService: moodleService}
}

func (s *MoodleController) CheckStatus(c *gin.Context) {
	result, err := s.moodleService.CheckStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, web.ApiResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
			Data:    nil,
		})
	}

	c.JSON(http.StatusOK, web.ApiResponse{
		Code:    "OK",
		Message: "OK",
		Data:    result,
	})
}

func (s *MoodleController) CreateUser(c *gin.Context) {
	var req web.MoodleUserCreateRequest

	fmt.Println("[DEBUG] Received Body:", req) // log

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateUser] Error: %v", err) // log
		c.JSON(http.StatusBadRequest, web.ApiResponse{
			Code:    "INVALID_PARAMS",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	log.Printf("[CreateUser] Received Request: %+v", req) // log

	// Call service
	result, err := s.moodleService.CreateUser(req)
	if err != nil {
		log.Println("[CreateUser] Error:", err)
		c.JSON(http.StatusInternalServerError, web.ApiResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, web.ApiResponse{
		Code:    "OK",
		Message: "OK",
		Data:    result,
	})
}

func (s *MoodleController) GetUserByField(c *gin.Context) {
	var req web.MoodleUserGetByFieldRequest

	// Bind JSON request body ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, web.ApiResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Format body request tidak valid.",
			Data:    err.Error(),
		})
		return
	}

	// service
	users, err := s.moodleService.GetUserByField(req)
	if err != nil {
		log.Printf("[DIAGNOSA] Controller menerima error. Tipe: %T, Isi: %v", err, err)
		if errors.Is(err, validation.ErrNotFound) {
			c.JSON(http.StatusNotFound, web.ApiResponse{
				Code:    "DATA_NOT_FOUND",
				Message: err.Error(), // Menggunakan pesan dari variabel ErrNotFound
			})
			return
		}

		if moodleErr, ok := err.(*web.MoodleException); ok {
			c.JSON(http.StatusBadRequest, web.ApiResponse{
				Code:    moodleErr.ErrorCode,
				Message: moodleErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, web.ApiResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Terjadi kesalahan pada server.",
		})
		return
	}

	c.JSON(http.StatusOK, web.ApiResponse{
		Code:    "OK",
		Message: "OK",
		Data:    users,
	})
}

func (s *MoodleController) UpdateUser(c *gin.Context) {
	var req []web.MoodleUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, web.ApiResponse{
			Code:    "INVALID_PARAMS",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	//periksa jika ada data
	err := s.moodleService.UpdateUsers(req)
	if err != nil {
		if moodleErr, ok := err.(*web.MoodleException); ok {
			c.JSON(http.StatusBadRequest, web.ApiResponse{
				Code:    moodleErr.ErrorCode,
				Message: moodleErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, web.ApiResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An internal error occurred",
			Data:    err.Error(), // Kirim pesan error internal untuk debug
		})
		return
	}

	c.JSON(http.StatusOK, web.ApiResponse{
		Code:    "OK",
		Message: "Users updated successfully",
	})
}

func (s *MoodleController) UserSync(c *gin.Context) {
	var req web.MoodleUserSyncRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, web.ApiResponse{
			Code:    "INVALID_REQUEST",
			Message: "Data yang dikirim tidak valid: " + err.Error(),
		})
		return
	}

	// Panggil service
	err := s.moodleService.UserSync(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, web.ApiResponse{
			Code:    "USER_SYNC_FAILED",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Berhasil
	c.JSON(http.StatusOK, web.ApiResponse{
		Code:    "OK",
		Message: "User synced successfully",
		Data:    nil,
	})
}

func (s *MoodleController) AssignRole(c *gin.Context) {
	var req web.MoodleRoleAssignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, web.ApiResponse{
			Code:    "INVALID_REQUEST",
			Message: "Data yang dikirim tidak valid: " + err.Error(),
		})
		return
	}

	log.Printf("[DEBUG] AssignRole Controller: Menerima request %+v", req)

	// panggil service
	if err := s.moodleService.AssignRole(req); err != nil {
		c.JSON(http.StatusBadRequest, web.ApiResponse{
			Code:    "USER_ASSIGN_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, web.ApiResponse{
		Code:    "OK",
		Message: "User assigned successfully",
	})
}
