import "./WebhookForm.scss"

import React, { useState } from "react"
import { Button, Field, Form, Input, Loader, Message, Select, SelectOption, TextArea, Toggle } from "@fider/components"
import { actions, Failure } from "@fider/services"
import { HStack, VStack } from "@fider/components/layout"
import { Webhook, WebhookPreviewResult, WebhookStatus, WebhookType } from "@fider/models"
import { HoverInfo } from "@fider/components/common/HoverInfo"
import { WebhookTemplateInfo } from "@fider/pages/Administration/components/webhook/WebhookTemplateInfo"

interface WebhookFormProps {
  webhook?: Webhook
  onSave: (data: WebhookFormState) => Promise<Failure | undefined>
  onCancel: () => void
}

export interface WebhookFormState extends Omit<Webhook, "id"> {
  error?: Failure
  typing?: NodeJS.Timeout
  preview?: WebhookPreviewResult | null
  isModalOpen: boolean
}

interface HttpHeaderProps {
  header?: string
  value?: string
  onEdit: (header: string, value: string) => void
  onRemove?: (header: string) => void
  allHeaders?: string[]
}

const HttpHeader = (props: HttpHeaderProps) => {
  const [header, setHeader] = useState(props.header || "")
  const [value, setValue] = useState(props.value || "")
  const [editing, setEditing] = useState(props.onRemove === undefined)

  const duplicate = props.allHeaders && props.allHeaders.includes(header)

  const suffix = editing ? (
    props.onRemove === undefined ? (
      <Button
        variant="primary"
        onClick={() => {
          props.onEdit(header, value)
          setHeader("")
          setValue("")
        }}
        disabled={duplicate || header.length === 0 || value.length === 0}
      >
        Add
      </Button>
    ) : (
      <Button
        variant="primary"
        onClick={() => {
          setEditing(false)
          props.onEdit(header, value)
        }}
        disabled={value.length === 0}
      >
        Save
      </Button>
    )
  ) : (
    <>
      <Button variant="secondary" onClick={() => setEditing(true)}>
        Edit
      </Button>
      <Button variant="danger" onClick={() => props.onRemove && props.onRemove(header)}>
        Remove
      </Button>
    </>
  )

  return (
    <div className="grid gap-4 grid-cols-2">
      <Input
        field={`header-${header}`}
        value={header}
        onChange={setHeader}
        placeholder="Header"
        disabled={props.onRemove !== undefined}
        suffix={duplicate ? "Duplicate" : undefined}
      />
      <Input field={`value-${header}`} value={value} onChange={setValue} placeholder="Value" disabled={!editing} suffix={suffix} />
    </div>
  )
}

export class WebhookForm extends React.Component<WebhookFormProps, WebhookFormState> {
  constructor(props: WebhookFormProps) {
    super(props)
    this.state = {
      name: props.webhook?.name || "",
      type: props.webhook?.type || WebhookType.NEW_POST,
      status: props.webhook?.status || WebhookStatus.DISABLED,
      url: props.webhook?.url || "",
      content: props.webhook?.content || "",
      http_method: props.webhook?.http_method || "POST",
      additional_http_headers: props.webhook?.additional_http_headers || {},
      isModalOpen: false,
    }
    this.calculatePreview().then((preview) => this.setState({ preview }))
  }

  public componentDidUpdate = (prevProps: WebhookFormProps, prevState: WebhookFormState) => {
    if (this.state.url !== prevState.url || this.state.content !== prevState.content) {
      if (this.state.typing) clearTimeout(this.state.typing)
      this.setTyping()
    }
  }

  private setTyping = () => {
    this.setState({
      preview: undefined,
      typing: setTimeout(async () => {
        this.setState({
          preview: await this.calculatePreview(),
          typing: undefined,
        })
      }, 2500),
    })
  }

  private calculatePreview = async () => {
    return actions.previewWebhook(this.state.type, this.state.url, this.state.content).then(
      (result) => (result.ok ? result.data : null),
      () => null
    )
  }

  private handleSave = async () => {
    const error = await this.props.onSave(this.state)
    if (error) {
      this.setState({ error })
    }
  }

  private handleCancel = async () => {
    this.props.onCancel()
  }

  private setName = (name: string) => {
    this.setState({ name })
  }

  private setType = (option?: SelectOption) => {
    this.setState({ type: option?.value as WebhookType })
  }

  private setStatus = (active: boolean) => {
    this.setState({ status: active ? WebhookStatus.ENABLED : WebhookStatus.DISABLED })
  }

  private setURL = (url: string) => {
    this.setState({ url })
  }

  private setContent = (content: string) => {
    this.setState({ content })
  }

  private setHttpMethod = (http_method: string) => {
    this.setState({ http_method })
  }

  private setHttpHeader = (header: string, value: string) => {
    this.setState((state) => ({
      additional_http_headers: {
        ...state.additional_http_headers,
        [header]: value,
      },
    }))
  }

  private removeHttpHeader = (header: string) => {
    this.setState((state) => {
      const { [header]: _, ...headers } = state.additional_http_headers
      return { additional_http_headers: headers }
    })
  }

  private showModal = () => {
    this.setState({ isModalOpen: true })
  }

  private hideModal = () => {
    this.setState({ isModalOpen: false })
  }

  public render() {
    const allHeaders = Object.keys(this.state.additional_http_headers)
    const title = this.props.webhook ? `Webhook #${this.props.webhook.id}: ${this.props.webhook.name}` : "New webhook"
    return (
      <>
        {this.state.status === WebhookStatus.FAILED && (
          <Message type="error" showIcon>
            This webhook has failed
          </Message>
        )}
        <h2 className="text-title mb-2">{title}</h2>
        <Form className="c-webhook-form" error={this.state.error}>
          <Input field="name" label="Name" value={this.state.name} onChange={this.setName} placeholder="My awesome webhook" />
          <Select
            label="Type"
            field="type"
            defaultValue={this.state.type}
            options={[
              { label: "New Post", value: WebhookType.NEW_POST },
              { label: "New Comment", value: WebhookType.NEW_COMMENT },
              { label: "Change Status", value: WebhookType.CHANGE_STATUS },
              { label: "Delete Post", value: WebhookType.DELETE_POST },
            ]}
            onChange={this.setType}
          />
          <Field label="Status">
            <Toggle active={this.state.status === WebhookStatus.ENABLED} onToggle={this.setStatus} />
            {this.state.status === WebhookStatus.FAILED && <p className="text-muted mt-1">This webhook was disabled due to a trigger failure</p>}
          </Field>
          <Input
            field="url"
            label="URL"
            afterLabel={<HoverInfo text="You can use Go template formatting with many properties here" onClick={this.showModal} />}
            value={this.state.url}
            onChange={this.setURL}
            placeholder="https://webhook.site/..."
          />
          <TextArea
            className="c-webhook-form__content"
            field="content"
            label="Content"
            afterLabel={<HoverInfo text="You can use Go template formatting with many properties here" onClick={this.showModal} />}
            value={this.state.content}
            onChange={this.setContent}
            placeholder="Request body"
          />
          <Input field="http_method" label="HTTP Method" value={this.state.http_method} onChange={this.setHttpMethod} placeholder="POST" />
          <Field label="Additional HTTP Headers" afterLabel={<HoverInfo text="Those headers are sent in the request when the webhook is triggered" />}>
            <VStack spacing={2}>
              {Object.entries(this.state.additional_http_headers).map(([header, value]) => (
                <HttpHeader key={header} header={header} value={value} onEdit={this.setHttpHeader} onRemove={this.removeHttpHeader} />
              ))}
              <HttpHeader onEdit={this.setHttpHeader} allHeaders={allHeaders} />
            </VStack>
          </Field>
          {(this.state.url || this.state.content) && (
            <Field label="Preview" className="c-webhook-form__preview">
              {this.state.preview === null ? (
                <p className="text-muted">Failed to load preview</p>
              ) : this.state.preview === undefined ? (
                <Loader className="text-center" text="Loading preview" />
              ) : (
                <VStack className="bg-gray-50 rounded-md p-2" spacing={2}>
                  {this.state.url && (
                    <div>
                      <h3 className="text-bold mb-1">URL</h3>
                      <pre>{this.state.preview.url.value ? this.state.preview.url.value : this.state.preview.url.error}</pre>
                      {this.state.preview.url.message && <p className="text-muted">{this.state.preview.url.message}</p>}
                    </div>
                  )}
                  {this.state.content && (
                    <div>
                      <h3 className="text-bold mb-1">Content</h3>
                      <pre>{this.state.preview.content.value ? this.state.preview.content.value : this.state.preview.content.error}</pre>
                      {this.state.preview.content.message && <p className="text-muted">{this.state.preview.content.message}</p>}
                    </div>
                  )}
                </VStack>
              )}
            </Field>
          )}
          <HStack>
            <Button variant="primary" onClick={this.handleSave}>
              Save
            </Button>
            <Button onClick={this.handleCancel} variant="tertiary">
              Cancel
            </Button>
          </HStack>
        </Form>
        <WebhookTemplateInfo type={this.state.type} isModalOpen={this.state.isModalOpen} onModalClose={this.hideModal} />
      </>
    )
  }
}
