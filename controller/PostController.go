package controller

import (
	"ginEssential/lxz/common"
	"ginEssential/lxz/model"
	"ginEssential/lxz/response"
	"ginEssential/lxz/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")

	// 创建post
	post := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	// 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
	}

	// 成功
	response.Success(ctx, nil, "创建成功")
}

func (p PostController) Update(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post

	if err := p.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, nil, "文章属于您，请勿非法操作")
		return
	}

	// 更新文章
	if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "更新成功")
}

func (p PostController) Show(ctx *gin.Context) {
	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	// if err := p.DB.Where("id = ?", postId).First(&post).Error; err != nil {
	if p.DB.Preload("Category").Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "成功")
}

func (p PostController) Delete(ctx *gin.Context) {
	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if err := p.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, nil, "文章属于您，请勿非法操作")
		return
	}

	p.DB.Delete(&post)

	response.Success(ctx, gin.H{"post": post}, "删除成功")
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 分页
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// 前端渲染分页需要知道总数
	var total int
	p.DB.Model(model.Post{}).Count(&total)

	// 返回数据
	response.Success(ctx, gin.H{"data": posts, "total": total}, "成功")
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{DB: db}
}
