import "./Loader.scss"

import React, { useState } from "react"
import { useTimeout } from "@fider/hooks"

interface LoaderProps {
  text?: string
  className?: string
}

export function Loader(props: LoaderProps) {
  const [show, setShow] = useState(false)

  useTimeout(() => {
    setShow(true)
  }, 500)

  return show ? (
    <div className={props.className}>
      <div className="c-loader" />
      {props.text && <span className="text-muted">{props.text}</span>}
    </div>
  ) : null
}
