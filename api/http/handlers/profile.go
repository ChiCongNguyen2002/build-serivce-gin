package handlers

import (
	"build-service-gin/api/http/models"
	"build-service-gin/common/custom/binding"
	"build-service-gin/internal/services"
	"build-service-gin/pkg/helpers/adapters"
	"build-service-gin/pkg/helpers/resp"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ProfileHandler struct {
	profileService services.IProfileService
}

func NewProfileHandler(profileService services.IProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (h *ProfileHandler) GetUserTransactionHistory(c *gin.Context) {
	var req models.GetUserTransactionHistoryReq

	if err := binding.GetBinding().Bind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		return
	}

	ctx := c.Request.Context()
	dataDomain := adapters.AdapterProfile{}.ConvReq2ServUserTransactionHistoryTx(req)
	data, total, err := h.profileService.GetUserHistoryByProfile(ctx, *dataDomain)
	if err != nil {
		c.JSON(http.StatusNotFound, resp.BuildErrorResp(err.ErrorCode, err.Description, resp.LangEN))
		return
	}

	rs := resp.BuildSuccessResp(resp.LangEN, data)
	rs.Paging = &resp.Paging{
		Total:  total,
		Offset: req.Offset,
		Limit:  req.Limit,
	}
	c.JSON(http.StatusOK, rs)
}

func (h *ProfileHandler) CreateUserTransactionHistory(c *gin.Context) {
	var req models.UserTransactionHistory
	if err := binding.GetBinding().Bind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		return
	}

	ctx := c.Request.Context()
	dataDomain := adapters.AdapterProfile{}.ConvModelToDomainUserTransactionHistoryTx(req)
	data, err := h.profileService.CreateUserTransactionHistory(ctx, dataDomain)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create user transaction history")
		c.JSON(http.StatusNotFound, resp.BuildErrorResp(err.ErrorCode, err.Description, resp.LangEN))
		return
	}

	successResp := resp.BuildSuccessResp(resp.LangEN, *data)
	c.JSON(http.StatusOK, successResp)
}

func (h *ProfileHandler) UpdateUserTransactionHistory(c *gin.Context) {
	var req models.UserTransactionHistory
	if err := binding.GetBinding().Bind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		return
	}

	ctx := c.Request.Context()
	dataDomain := adapters.AdapterProfile{}.ConvModelToDomainUserTransactionHistoryTx(req)
	data, err := h.profileService.UpdateUserTransactionHistoryByProfile(ctx, dataDomain, dataDomain.ProfileID)
	if err != nil {
		c.JSON(http.StatusNotFound, resp.BuildErrorResp(err.ErrorCode, err.Description, resp.LangEN))
		return
	}
	c.JSON(http.StatusOK, resp.BuildSuccessResp(resp.LangEN, data))
}

func (h *ProfileHandler) DeleteUserTransactionHistory(c *gin.Context) {
	var req models.GetUserTransactionHistoryByProfileReq
	if err := binding.GetBinding().Bind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		return
	}

	ctx := c.Request.Context()
	err := h.profileService.DeleteUserTransactionHistoryByProfile(ctx, req.ProfileID)
	if err != nil {
		c.JSON(http.StatusNotFound, resp.BuildErrorResp(err.ErrorCode, err.Description, resp.LangEN))
		return
	}
	c.JSON(http.StatusOK, resp.BuildSuccessResp(resp.LangEN, nil))
}

func (h *ProfileHandler) GetUserTransactionHistoryPostgres(c *gin.Context) {
	var req models.GetUserTransactionHistoryReq

	if err := binding.GetBinding().Bind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		return
	}

	ctx := c.Request.Context()
	dataDomain := adapters.AdapterProfile{}.ConvReq2ServUserTransactionHistoryTx(req)
	data, total, err := h.profileService.GetUserHistoryByProfilePostgresql(ctx, *dataDomain)
	if err != nil {
		c.JSON(http.StatusNotFound, resp.BuildErrorResp(err.ErrorCode, err.Description, resp.LangEN))
		return
	}

	rs := resp.BuildSuccessResp(resp.LangEN, data)
	rs.Paging = &resp.Paging{
		Total:  total,
		Offset: req.Offset,
		Limit:  req.Limit,
	}
	c.JSON(http.StatusOK, rs)
}

func (h *ProfileHandler) CreateUserTransactionHistoryPostgres(c *gin.Context) {
	var req models.UserTransactionHistory
	if err := binding.GetBinding().Bind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BuildErrorResp(resp.ErrDataInvalid, err.Error(), resp.LangEN))
		return
	}

	ctx := c.Request.Context()
	dataDomain := adapters.AdapterProfile{}.ConvModelToDomainUserTransactionHistoryTx(req)
	data, err := h.profileService.CreateUserTransactionHistoryPostgresql(ctx, dataDomain)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create user transaction history")
		c.JSON(http.StatusNotFound, resp.BuildErrorResp(err.ErrorCode, err.Description, resp.LangEN))
		return
	}

	successResp := resp.BuildSuccessResp(resp.LangEN, *data)
	c.JSON(http.StatusOK, successResp)
}
