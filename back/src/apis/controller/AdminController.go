package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"../auth"
	"../db"
	"../model"
)

type AdminController struct{}

type ProgressResponse struct {
	AllProgress []ProgressUserId `json:"prog"`
}

type ProgressUserId struct {
	UserId     int   `json:"uid"`
	QuestionId []int `json:"progress"`
}

type QuestionHints struct {
	Hint string `json:"hint"`
}

type QuestionRequest struct {
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Answer  string         `json:"answer"`
	Hints   []QuestionHint `json:"hints"`
}

func (pc AdminController) AddQuestion(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	_, err := auth.VerifyToken(tokenString)
	if err != nil {
		c.String(http.StatusUnauthorized, "token is invalid")
		return
	}

	db := db.GormConnect()
	request := QuestionRequest{}
	err = c.BindJSON(&request)
	if err != nil {
		c.String(http.StatusBadRequest, "request failed(json error)")
	}
	problem := model.Problem{}
	problem.Pro_Title = request.Title
	problem.Pro_Content = request.Content
	problem.Pro_Answer = request.Answer
	db.NewRecord(problem)
	db.Create(&problem)
	db.First(&problem)
	hint := model.Hint{}
	for _, h := range request.Hints {
		hint.ID = 0
		hint.Hint = h.Hint
		hint.Pro_Id = int(problem.ID)
		db.NewRecord(hint)
		db.Create(&hint)
	}
	c.String(http.StatusCreated, "add complete")
}

func (pc AdminController) Delete(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	_, err := auth.VerifyToken(tokenString)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	problem := model.Problem{}
	u64, err := strconv.ParseUint(c.Param("id"), 10, 0)
	problem.ID = uint(u64)

	db := db.GormConnect()
	db.Delete(&problem)
	c.String(http.StatusCreated, "complete delete")
}

func (pc AdminController) EditQuestion(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	_, err := auth.VerifyToken(tokenString)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	edit := QuestionRequest{}
	err = c.BindJSON(&edit)
	if err != nil {
		c.String(http.StatusBadRequest, "request failed(json error)")
	}

	db := db.GormConnect()
	problem := model.Problem{}
	u64, err := strconv.ParseUint(c.Param("id"), 10, 0)
	problem.ID = uint(u64)

	problem_after := problem
	db.First(&problem_after)
	problem_after.Pro_Title = edit.Title
	problem_after.Pro_Content = edit.Content
	problem_after.Pro_Answer = edit.Answer
	db.Model(&problem).Update(problem_after)
	db.Save(&problem_after)
	c.String(http.StatusCreated, "complete edit")

	hint := model.Hint{}
	db.Where("pro_id=?", problem.ID).Delete(&hint)
	for _, h := range edit.Hints {
		hint_after := model.Hint{}
		hint_after.Hint = h.Hint
		hint_after.Pro_Id = int(problem.ID)
		db.Create(&hint_after)
	}
	c.String(http.StatusCreated, "complete edit")

}

func (pc AdminController) AllProgress(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	_, err := auth.VerifyToken(tokenString)
	if err != nil {
		c.String(http.StatusUnauthorized, "token is invalid")
		return
	}

	db := db.GormConnect()
	user := []model.User{}
	prog := []ProgressUserId{}
	var count int
	db.Model(&user).Count(&count)
	db.Find(&user)

	for _, h := range user {
		progress := []model.Progress{}
		var cnt int
		db.Model(&progress).Where("user_id=?", h.ID).Count(&cnt)
		db.Where("user_id=?", h.ID).Find(&progress)
		var proid []int
		for j := 0; j < cnt; j++ {
			proid = append(proid, progress[j].Pro_Id)
		}

		prog = append(prog, ProgressUserId{
			UserId:     progress[0].User_Id,
			QuestionId: proid,
		})
	}
	adf := ProgressResponse{
		prog,
	}

	c.JSON(http.StatusOK, adf)
}

func (pc AdminController) UserIdProgress(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	_, err := auth.VerifyToken(tokenString)
	if err != nil {
		c.String(http.StatusUnauthorized, "token is invalid")
		return
	}

	u64, err := strconv.ParseUint(c.Param("id"), 10, 0)
	userid := uint(u64)
	db := db.GormConnect()
	progress := []model.Progress{}
	db.Find(&progress, "user_id=?", userid)
	var count int
	db.Model(&progress).Where("user_id = ?", userid).Count(&count)

	var proid []int
	for i := 0; i < count; i++ {
		proid = append(proid, progress[i].Pro_Id)
	}

	type ResponseProgress struct {
		QuestionId []int `json:"progress"`
	}

	adf := ResponseProgress{
		proid,
	}

	c.JSON(http.StatusOK, adf)
}
