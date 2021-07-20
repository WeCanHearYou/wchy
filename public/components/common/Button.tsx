import "./Button.scss"

import React from "react"
import { classSet } from "@fider/services"

interface ButtonProps {
  className?: string
  disabled?: boolean
  href?: string
  rel?: "nofollow"
  target?: "_self" | "_blank" | "_parent" | "_top"
  type?: "button" | "submit"
  variant?: "primary" | "danger" | "secondary" | "tertiary"
  size?: "small" | "default" | "large"
  onClick?: (event: ButtonClickEvent) => Promise<any> | void
}

interface ButtonState {
  clicked: boolean
}

export class ButtonClickEvent {
  private shouldEnable = true
  public preventEnable(): void {
    this.shouldEnable = false
  }
  public canEnable(): boolean {
    return this.shouldEnable
  }
}

export class Button extends React.Component<ButtonProps, ButtonState> {
  private unmounted = false

  public static defaultProps: Partial<ButtonProps> = {
    size: "default",
    variant: "secondary",
    type: "button",
  }

  public constructor(props: ButtonProps) {
    super(props)
    this.state = {
      clicked: false,
    }
  }

  public componentWillUnmount() {
    this.unmounted = true
  }

  public click = async (e?: React.SyntheticEvent<HTMLElement>) => {
    if (e) {
      e.preventDefault()
      e.stopPropagation()
    }

    if (this.state.clicked) {
      return
    }

    const event = new ButtonClickEvent()
    this.setState({ clicked: true })
    if (this.props.onClick) {
      await this.props.onClick(event)
      if (!this.unmounted && event.canEnable()) {
        this.setState({ clicked: false })
      }
    }
  }

  public render() {
    const className = classSet({
      "c-button": true,
      [`c-button--${this.props.size}`]: this.props.size,
      [`c-button--${this.props.variant}`]: this.props.variant,
      "c-button--loading": this.state.clicked,
      "c-button--disabled": this.state.clicked || this.props.disabled,
      [this.props.className || ""]: this.props.className,
      "shadow-sm": this.props.variant !== "tertiary",
    })

    if (this.props.href) {
      return (
        <a href={this.props.href} rel={this.props.rel} target={this.props.target} className={className}>
          {this.props.children}
        </a>
      )
    } else if (this.props.onClick) {
      return (
        <button type={this.props.type} className={className} onClick={this.click}>
          {this.props.children}
        </button>
      )
    } else {
      return (
        <button type={this.props.type} className={className}>
          {this.props.children}
        </button>
      )
    }
  }
}
