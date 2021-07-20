import "./WebhookProperties.scss"

import React from "react"

import { VStack } from "@fider/components/layout"
import { StringObject } from "@fider/services"

interface WebhookPropertiesProps {
  properties: StringObject
  propsName: string
  valueName: string
}

interface PropertyProps {
  value: any
}

const Property = (props: PropertyProps) => {
  const grayValText = (txt: string) => <span className="c-webhook-properties__val--gray c-webhook-properties__val--italic">&lt;{txt}&gt;</span>

  if (Array.isArray(props.value))
    return (
      <VStack spacing={2} divide>
        {props.value.map((val, i) => (
          <Property key={i} value={val} />
        ))}
      </VStack>
    )

  if (props.value === "") return grayValText("empty")
  if (props.value === null) return grayValText("null")
  if (props.value === undefined) return grayValText("undefined")
  if (props.value === true) return <span className="c-webhook-properties__val--green">true</span>
  if (props.value === false) return <span className="c-webhook-properties__val--red">false</span>

  const type = typeof props.value
  switch (type) {
    case "string":
      return <span>{props.value}</span>
    case "number":
    case "bigint":
      return <span className="c-webhook-properties__val--blue">{props.value}</span>
    case "object":
      return <WebhookProperties properties={props.value} propsName="Name" valueName="Value" />
    default:
      return props.value
  }
}

export const WebhookProperties = (props: WebhookPropertiesProps) => {
  return (
    <table className="c-webhook-properties">
      <thead>
        <tr>
          <th>{props.propsName}</th>
          <th>{props.valueName}</th>
        </tr>
      </thead>
      <tbody>
        {Object.entries(props.properties).map(([prop, val]) => (
          <tr key={prop}>
            <td className="c-webhook-properties__prop">{prop}</td>
            <td className="c-webhook-properties__val">
              <Property value={val} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
