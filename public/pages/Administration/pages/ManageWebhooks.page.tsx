import React from "react"
import { Button } from "@fider/components"

import { Webhook, WebhookStatus, WebhookType } from "@fider/models"
import { actions, Failure } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"
import { WebhookForm, WebhookFormState } from "../components/webhook/WebhookForm"
import { WebhookListItem } from "../components/webhook/WebhookListItem"
import { VStack } from "@fider/components/layout"

interface ManageWebhooksPageProps {
  webhooks: Webhook[]
}

interface ManageWebhooksPageState {
  isAdding: boolean
  allWebhooks: Webhook[]
  editing?: Webhook
}

const webhookSorter = (w1: Webhook, w2: Webhook) => {
  if (w1.name < w2.name) {
    return -1
  } else if (w1.name > w2.name) {
    return 1
  }
  return 0
}

interface WebhooksListProps {
  title: string
  description: string
  list: JSX.Element[]
}

const WebhooksList = (props: WebhooksListProps) => {
  return (
    <div>
      <h2 className="text-display">{props.title}</h2>
      <p className="text-muted">These webhooks are triggered every time {props.description}.</p>
      <VStack spacing={4} divide={true}>
        {props.list.length === 0 ? <p className="text-muted">There aren’t any &quot;{props.title.toLowerCase()}&quot; webhook yet.</p> : props.list}
      </VStack>
    </div>
  )
}

export default class ManageWebhooksPage extends AdminBasePage<ManageWebhooksPageProps, ManageWebhooksPageState> {
  public id = "p-admin-webhooks"
  public name = "webhooks"
  public title = "Webhooks"
  public subtitle = "Manage your site webhooks"

  constructor(props: ManageWebhooksPageProps) {
    super(props)
    this.state = {
      isAdding: false,
      allWebhooks: this.props.webhooks.sort(webhookSorter),
    }
  }

  private addNew = async () => {
    this.setState({
      isAdding: true,
      editing: undefined,
    })
  }

  private cancelAdd = () => {
    this.setState({ isAdding: false })
  }

  private saveNewWebhook = async (data: WebhookFormState): Promise<Failure | undefined> => {
    const result = await actions.createWebhook(data.name, data.type, data.status, data.url, data.content, data.http_method, data.additional_http_headers)
    if (result.ok) {
      this.setState({
        isAdding: false,
        allWebhooks: this.state.allWebhooks.concat({ id: result.data.id, ...data }).sort(webhookSorter),
      })
    } else {
      return result.error
    }
  }

  private startWebhookEditing = (webhook: Webhook) => {
    this.setState({ editing: webhook })
  }

  private cancelEdit = () => {
    this.setState({ editing: undefined })
  }

  private handleWebhookDeleted = (webhook: Webhook) => {
    const idx = this.state.allWebhooks.indexOf(webhook)
    this.setState({
      allWebhooks: this.state.allWebhooks.splice(idx, 1) && this.state.allWebhooks,
    })
  }

  private handleWebhookEdited = async (data: WebhookFormState): Promise<Failure | undefined> => {
    const webhook = this.state.editing
    if (webhook === undefined) return // impossible
    const result = await actions.updateWebhook(
      webhook.id,
      data.name,
      data.type,
      data.status,
      data.url,
      data.content,
      data.http_method,
      data.additional_http_headers
    )
    if (result.ok) {
      webhook.name = data.name
      webhook.type = data.type
      webhook.status = data.status === WebhookStatus.FAILED ? WebhookStatus.DISABLED : data.status
      webhook.url = data.url
      webhook.content = data.content
      webhook.http_method = data.http_method
      webhook.additional_http_headers = data.additional_http_headers

      this.setState({
        editing: undefined,
        allWebhooks: this.state.allWebhooks.sort(webhookSorter),
      })
    } else {
      return result.error
    }
  }

  private handleWebhookFailed = (webhook: Webhook) => {
    webhook.status = WebhookStatus.FAILED
    this.setState({
      allWebhooks: this.state.allWebhooks.sort(webhookSorter),
    })
  }

  private getWebhookList(filter: (webhook: Webhook) => boolean) {
    return this.state.allWebhooks.filter(filter).map((w) => {
      return (
        <WebhookListItem
          key={w.id}
          webhook={w}
          onWebhookDeleted={this.handleWebhookDeleted}
          editWebhook={this.startWebhookEditing}
          onWebhookFailed={this.handleWebhookFailed}
        />
      )
    })
  }

  public content() {
    if (this.state.isAdding) {
      return <WebhookForm onSave={this.saveNewWebhook} onCancel={this.cancelAdd} />
    }

    if (this.state.editing) {
      return <WebhookForm onSave={this.handleWebhookEdited} onCancel={this.cancelEdit} webhook={this.state.editing} />
    }

    const newPostList = this.getWebhookList((w) => w.type === WebhookType.NEW_POST)
    const newCommentList = this.getWebhookList((w) => w.type === WebhookType.NEW_COMMENT)
    const changeStatusList = this.getWebhookList((w) => w.type === WebhookType.CHANGE_STATUS)
    const deletePostList = this.getWebhookList((w) => w.type === WebhookType.DELETE_POST)

    return (
      <VStack spacing={8}>
        <WebhooksList title="New Post" description="a new post is created on this site" list={newPostList} />
        <WebhooksList title="New Comment" description="a new comment is created on any post" list={newCommentList} />
        <WebhooksList title="Change Status" description="the status of a post is changed" list={changeStatusList} />
        <WebhooksList title="Delete Post" description="a post is deleted on this site" list={deletePostList} />
        <div>
          <Button variant="secondary" onClick={this.addNew}>
            Add new
          </Button>
        </div>
      </VStack>
    )
  }
}
