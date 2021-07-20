import React from "react"
import { classSet } from "@fider/services"
import { ValidationContext } from "../"
import { DisplayError, hasError } from "./DisplayError"
import Textarea from "react-textarea-autosize"

import "./TextArea.scss"

interface TextAreaProps {
  label?: string
  field: string
  value?: string
  disabled?: boolean
  minRows?: number
  placeholder?: string
  afterLabel?: JSX.Element
  onChange?: (value: string) => void
  inputRef?: React.MutableRefObject<any>
  onFocus?: React.FocusEventHandler<HTMLTextAreaElement>
  className?: string
}

export const TextArea: React.FunctionComponent<TextAreaProps> = (props) => {
  const onChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    if (props.onChange) {
      props.onChange(e.currentTarget.value)
    }
  }

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <>
          <div className="c-form-field">
            {!!props.label && (
              <label htmlFor={`input-${props.field}`}>
                {props.label}
                {props.afterLabel}
              </label>
            )}
            <Textarea
              className={classSet({
                "c-textarea": true,
                "c-textarea--error": hasError(props.field, ctx.error),
                [props.className || ""]: props.className,
              })}
              id={`input-${props.field}`}
              disabled={props.disabled}
              onChange={onChange}
              value={props.value}
              minRows={props.minRows || 3}
              placeholder={props.placeholder}
              ref={props.inputRef}
              onFocus={props.onFocus}
            />
            <DisplayError fields={[props.field]} error={ctx.error} />
            {props.children}
          </div>
        </>
      )}
    </ValidationContext.Consumer>
  )
}
