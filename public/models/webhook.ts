import { StringObject } from "@fider/services"

export interface Webhook {
  id: number
  name: string
  type: WebhookType
  status: WebhookStatus
  url: string
  content: string
  http_method: string
  additional_http_headers: HttpHeaders
}

export enum WebhookType {
  NEW_POST = "new_post",
  NEW_COMMENT = "new_comment",
  CHANGE_STATUS = "change_status",
  DELETE_POST = "delete_post",
}

export enum WebhookStatus {
  ENABLED = "enabled",
  DISABLED = "disabled",
  FAILED = "failed",
}

export type HttpHeaders = StringObject<string>

export interface WebhookTriggerResult {
  webhook: Webhook
  props: StringObject
  success: boolean
  url: string
  content: string
  status_code: number
  message: string
  error: string
}

export interface WebhookPreviewResult {
  url: PreviewedField
  content: PreviewedField
}

export interface PreviewedField {
  value: string
  message: string
  error: string
}
