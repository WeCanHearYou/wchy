import { http, Result, StringObject } from "@fider/services"
import { HttpHeaders, WebhookPreviewResult, WebhookStatus, WebhookTriggerResult, WebhookType } from "@fider/models"

export const createWebhook = async (
  name: string,
  type: WebhookType,
  status: WebhookStatus,
  url: string,
  content: string,
  http_method: string,
  additional_http_headers: HttpHeaders
): Promise<Result<{ id: number }>> => {
  return await http.post(`/_api/admin/webhook`, { name, type, status, url, content, http_method, additional_http_headers })
}

export const updateWebhook = async (
  id: number,
  name: string,
  type: WebhookType,
  status: WebhookStatus,
  url: string,
  content: string,
  http_method: string,
  additional_http_headers: HttpHeaders
): Promise<Result> => {
  return await http.put(`/_api/admin/webhook/${id}`, { name, type, status, url, content, http_method, additional_http_headers })
}

export const deleteWebhook = async (id: number): Promise<Result> => {
  return await http.delete(`/_api/admin/webhook/${id}`)
}

export const triggerWebhook = async (id: number): Promise<Result<WebhookTriggerResult>> => {
  return await http.get(`/_api/admin/webhook/trigger/${id}`)
}

export const previewWebhook = async (type: WebhookType, url: string, content: string): Promise<Result<WebhookPreviewResult>> => {
  return await http.post("/_api/admin/webhook/preview", { type, url, content })
}

export const getWebhookHelp = async (type: WebhookType): Promise<Result<StringObject>> => {
  return await http.get(`/_api/admin/webhook/props/${type}`)
}
