package handlers

import (
	"build-service-gin/api/http/models"
	"build-service-gin/common/custom/binding"
	"build-service-gin/internal/services"
	"build-service-gin/pkg/helpers/adapters"
	"build-service-gin/pkg/helpers/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PointHandler struct {
	pointService services.IPointService
}

func NewPointHandler(pointService services.IPointService) *PointHandler {
	return &PointHandler{
		pointService: pointService,
	}
}

func (h *PointHandler) CreatePointTransaction(c *gin.Context) {
	var req *models.OrderRequest

	if err := binding.GetBinding().Bind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		return
	}

	dataDomain := adapters.AdapterLPPoint{}.ConvertOrderHandler2Domain(req)

	err := h.pointService.CreatePointTransaction(c.Request.Context(), dataDomain)
	if err != nil {
		c.JSON(http.StatusNotFound, resp.BuildErrorResp(err.ErrorCode, err.Description, resp.LangEN))
		return
	}

	c.JSON(http.StatusOK, resp.BuildSuccessResp(resp.LangEN, nil))
}
