import "./HoverInfo.scss"

import React from "react"
import { Icon } from "./Icon"

import IconInformationCircle from "@fider/assets/images/heroicons-information-circle.svg"

interface InfoProps {
  text: string
  onClick?: () => void
  href?: string
  target?: "_self" | "_blank" | "_parent" | "_top"
}

export const HoverInfo = (props: InfoProps) => {
  const Elem = props.href ? "a" : "span"
  return (
    <Elem className="c-hoverinfo" data-tooltip={props.text} onClick={props.onClick} href={props.href} target={props.target}>
      <Icon width="15" height="15" className="c-hoverinfo__icon" sprite={IconInformationCircle} />
    </Elem>
  )
}
