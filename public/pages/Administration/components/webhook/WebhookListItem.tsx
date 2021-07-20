import "./WebhookListItem.scss"

import React, { useState } from "react"
import { Webhook, WebhookStatus, WebhookTriggerResult } from "@fider/models"
import { Button, Icon } from "@fider/components"
import { actions, notify } from "@fider/services"

import IconX from "@fider/assets/images/heroicons-x.svg"
import IconPencilAlt from "@fider/assets/images/heroicons-pencil-alt.svg"
import IconPlay from "@fider/assets/images/heroicons-play.svg"
import IconCheckCircle from "@fider/assets/images/heroicons-check-circle.svg"
import IconXCircle from "@fider/assets/images/heroicons-x-circle.svg"
import IconExclamation from "@fider/assets/images/heroicons-exclamation.svg"
import { HStack, VStack } from "@fider/components/layout"
import { WebhookFailInfo } from "./WebhookFailInfo"

interface WebhookListItemProps {
  webhook: Webhook
  editWebhook: (webhook: Webhook) => void
  onWebhookDeleted: (webhook: Webhook) => void
  onWebhookFailed: (webhook: Webhook) => void
}

interface WebhookIconProps {
  status: WebhookStatus
}

const WebhookIcon = (props: WebhookIconProps) => {
  let text, icon
  switch (props.status) {
    case WebhookStatus.ENABLED:
      text = "Enabled"
      icon = IconCheckCircle
      break
    case WebhookStatus.DISABLED:
      text = "Disabled"
      icon = IconXCircle
      break
    case WebhookStatus.FAILED:
      text = "Failed"
      icon = IconExclamation
  }
  return (
    <div data-tooltip={text}>
      <Icon width="23" height="23" className={`c-webhook-listitem__icon c-webhook-listitem__icon--${text.toLowerCase()}`} sprite={icon} />
    </div>
  )
}

export const WebhookListItem = (props: WebhookListItemProps) => {
  const [deleting, setDeleting] = useState(false)
  const [triggerResult, setTriggerResult] = useState<WebhookTriggerResult | undefined>(undefined)
  const [isFailInfoModalOpen, setIsFailInfoModalOpen] = useState(false)

  const showFailInfoModal = () => setIsFailInfoModalOpen(true)
  const hideFailInfoModal = () => setIsFailInfoModalOpen(false)

  const deleteWebhook = async () => {
    const result = await actions.deleteWebhook(props.webhook.id)
    if (result.ok) {
      setDeleting(false)
      props.onWebhookDeleted(props.webhook)
    }
  }

  const triggerWebhook = async () => {
    const result = await actions.triggerWebhook(props.webhook.id)
    setTriggerResult(result.data)
    if (result.ok && result.data.success) {
      notify.success("Successfully triggered webhook")
    } else {
      notify.error(result.data.message)
      props.onWebhookFailed(props.webhook)
    }
  }

  const renderDeleteMode = () => {
    return (
      <VStack spacing={2}>
        <div>
          <b>Are you sure?</b>{" "}
          <span>
            The webhook #{props.webhook.id} &quot;{props.webhook.name}&quot; will be deleted forever. Alternatively, you may want to <b>disable</b> it instead.
          </span>
        </div>
        <div>
          <Button variant="danger" onClick={deleteWebhook}>
            Delete webhook
          </Button>
          <Button onClick={() => setDeleting(false)} variant="tertiary">
            Cancel
          </Button>
        </div>
      </VStack>
    )
  }

  const renderViewMode = () => {
    return (
      <HStack justify="between">
        <HStack>
          <WebhookIcon status={props.webhook.status} />
          <h3 className="text-title">
            <span className="text-muted">#{props.webhook.id}</span> {props.webhook.name}
          </h3>
          {triggerResult?.success === false && (
            <WebhookFailInfo result={triggerResult} isModalOpen={isFailInfoModalOpen} onModalOpen={showFailInfoModal} onModalClose={hideFailInfoModal} />
          )}
        </HStack>
        <HStack>
          <Button size="small" onClick={triggerWebhook}>
            <Icon sprite={IconPlay} />
            <span>Trigger</span>
          </Button>
          <Button size="small" onClick={() => props.editWebhook(props.webhook)}>
            <Icon sprite={IconPencilAlt} />
            <span>Edit</span>
          </Button>
          <Button size="small" onClick={() => setDeleting(true)}>
            <Icon sprite={IconX} />
            <span>Delete</span>
          </Button>
        </HStack>
      </HStack>
    )
  }

  return deleting ? renderDeleteMode() : renderViewMode()
}
