package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"

	"github.com/didi/nightingale/v5/src/models"
)

// Return all, front-end search and paging
func alertSubscribeGets(c *gin.Context) {
	bgid := ginx.UrlParamInt64(c, "id")
	lst, err := models.AlertSubscribeGets(bgid)
	if err == nil {
		ugcache := make(map[int64]*models.UserGroup)
		for i := 0; i < len(lst); i++ {
			lst[i].FillUserGroups(ugcache)
		}

		rulecache := make(map[int64]string)
		for i := 0; i < len(lst); i++ {
			lst[i].FillRuleName(rulecache)
		}
	}
	ginx.NewRender(c).Data(lst, err)
}

func alertSubscribeAdd(c *gin.Context) {
	var f models.AlertSubscribe
	ginx.BindJSON(c, &f)

	if f.GroupId <= 0 {
		ginx.Bomb(http.StatusBadRequest, "group_id invalid")
	}

	username := c.MustGet("username").(string)
	f.CreateBy = username
	f.UpdateBy = username
	f.GroupId = ginx.UrlParamInt64(c, "id")

	ginx.NewRender(c).Message(f.Add())
}

func alertSubscribePut(c *gin.Context) {
	var fs []models.AlertSubscribe
	ginx.BindJSON(c, &fs)

	timestamp := time.Now().Unix()
	username := c.MustGet("username").(string)
	for i := 0; i < len(fs); i++ {
		fs[i].UpdateBy = username
		fs[i].UpdateAt = timestamp
		ginx.Dangerous(fs[i].Update(
			"rule_id",
			"tags",
			"redefine_severity",
			"new_severity",
			"redefine_channels",
			"new_channels",
			"user_group_ids",
			"update_at",
			"update_by",
		))
	}

	ginx.NewRender(c).Message(nil)
}

func alertSubscribeDel(c *gin.Context) {
	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()

	ginx.NewRender(c).Message(models.AlertSubscribeDel(f.Ids))
}