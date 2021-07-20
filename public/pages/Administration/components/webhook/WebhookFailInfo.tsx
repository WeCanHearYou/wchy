import "./WebhookFailInfo.scss"

import React from "react"
import { WebhookTriggerResult } from "@fider/models"
import { Button, Modal } from "@fider/components"

import { VStack } from "@fider/components/layout"
import { HoverInfo } from "@fider/components/common/HoverInfo"
import { WebhookProperties } from "@fider/pages/Administration/components/webhook/WebhookProperties"

interface WebhookFailInfoProps {
  result: WebhookTriggerResult
  isModalOpen: boolean
  onModalOpen: () => void
  onModalClose: () => void
}

interface InfoPropertyProps {
  value: string | number
  name: string
  info: string
  multiline?: boolean
}

const InfoProperty = (props: InfoPropertyProps) => {
  return props.value ? (
    <div>
      <h3 className="text-title mb-1">
        {props.name}
        <HoverInfo text={props.info} />
      </h3>
      {props.multiline ? <pre>{props.value}</pre> : <p>{props.value}</p>}
    </div>
  ) : (
    <div className="text-muted">
      <span className="text-bold">{props.name}</span> info not available
    </div>
  )
}

export const WebhookFailInfo = (props: WebhookFailInfoProps) => {
  return (
    <>
      <HoverInfo text="Click to show failure details" onClick={props.onModalOpen} />
      <Modal.Window isOpen={props.isModalOpen} onClose={props.onModalClose} size="large">
        <Modal.Header>Webhook trigger failure details</Modal.Header>
        <Modal.Content>
          <VStack className="c-webhook-failinfo" spacing={4} divide>
            <InfoProperty value={props.result.message} name="Message" info="Generic information about where it failed" />
            <InfoProperty value={props.result.error} name="Error" info="Detailed information about what failed" multiline />
            <InfoProperty value={props.result.url} name="URL" info="Parsed URL where the request has been made" />
            <InfoProperty value={props.result.content} name="Content" info="Parsed content that was sent as request body" multiline />
            <InfoProperty value={props.result.status_code} name="Status code" info="HTTP response status code of the request" />
            <div>
              <h3 className="text-title mb-1">
                Properties
                <HoverInfo text="Properties used when parsing URL and content" />
              </h3>
              <WebhookProperties properties={props.result.props} propsName="Property name" valueName="Resolved value" />
            </div>
          </VStack>
        </Modal.Content>
        <Modal.Footer>
          <Button variant="tertiary" onClick={props.onModalClose}>
            Close
          </Button>
        </Modal.Footer>
      </Modal.Window>
    </>
  )
}
