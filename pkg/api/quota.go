package api

import (
	macaron "gopkg.in/macaron.v1"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/setting"
)

func (hs *HTTPServer) GetCurrentOrgQuotas(c *models.ReqContext) response.Response {
	return hs.getOrgQuotasHelper(c, c.OrgId)
}

func (hs *HTTPServer) GetOrgQuotas(c *models.ReqContext) response.Response {
	return hs.getOrgQuotasHelper(c, c.ParamsInt64(":orgId"))
}

func (hs *HTTPServer) getOrgQuotasHelper(c *models.ReqContext, orgID int64) response.Response {
	hasAccess := accesscontrol.HasAccess(hs.AccessControl, c)
	if !hasAccess(accesscontrol.NoReq, accesscontrol.EvalPermission(ActionOrgsRead, buildOrgsIdScope(c.OrgId))) {
		return response.Error(403, "Access denied to org", nil)
	}

	if !hs.Cfg.Quota.Enabled {
		return response.Error(404, "Quotas not enabled", nil)
	}
	query := models.GetOrgQuotasQuery{OrgId: orgID}

	if err := hs.SQLStore.GetOrgQuotas(c.Req.Context(), &query); err != nil {
		return response.Error(500, "Failed to get org quotas", err)
	}

	return response.JSON(200, query.Result)
}

func UpdateOrgQuota(c *models.ReqContext, cmd models.UpdateOrgQuotaCmd) response.Response {
	if !setting.Quota.Enabled {
		return response.Error(404, "Quotas not enabled", nil)
	}
	cmd.OrgId = c.ParamsInt64(":orgId")
	cmd.Target = macaron.Params(c.Req)[":target"]

	if _, ok := setting.Quota.Org.ToMap()[cmd.Target]; !ok {
		return response.Error(404, "Invalid quota target", nil)
	}

	if err := bus.DispatchCtx(c.Req.Context(), &cmd); err != nil {
		return response.Error(500, "Failed to update org quotas", err)
	}
	return response.Success("Organization quota updated")
}

func GetUserQuotas(c *models.ReqContext) response.Response {
	if !setting.Quota.Enabled {
		return response.Error(404, "Quotas not enabled", nil)
	}
	query := models.GetUserQuotasQuery{UserId: c.ParamsInt64(":id")}

	if err := bus.DispatchCtx(c.Req.Context(), &query); err != nil {
		return response.Error(500, "Failed to get org quotas", err)
	}

	return response.JSON(200, query.Result)
}

func UpdateUserQuota(c *models.ReqContext, cmd models.UpdateUserQuotaCmd) response.Response {
	if !setting.Quota.Enabled {
		return response.Error(404, "Quotas not enabled", nil)
	}
	cmd.UserId = c.ParamsInt64(":id")
	cmd.Target = macaron.Params(c.Req)[":target"]

	if _, ok := setting.Quota.User.ToMap()[cmd.Target]; !ok {
		return response.Error(404, "Invalid quota target", nil)
	}

	if err := bus.DispatchCtx(c.Req.Context(), &cmd); err != nil {
		return response.Error(500, "Failed to update org quotas", err)
	}
	return response.Success("Organization quota updated")
}
