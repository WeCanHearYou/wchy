import "./SignInControl.scss"

import React, { useState } from "react"
import { SocialSignInButton, Form, Button, Input, Message } from "@fider/components"
import { Divider } from "@fider/components/layout"
import { device, actions, Failure, isCookieEnabled } from "@fider/services"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/macro"

interface SignInControlProps {
  useEmail: boolean
  redirectTo?: string
  onEmailSent?: (email: string) => void
}

export const SignInControl: React.FunctionComponent<SignInControlProps> = (props) => {
  const fider = useFider()
  const [showEmailForm, setShowEmailForm] = useState(fider.session.tenant.isEmailAuthAllowed)
  const [email, setEmail] = useState("")
  const [error, setError] = useState<Failure | undefined>(undefined)

  const forceShowEmailForm = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault()
    setShowEmailForm(true)
  }

  const signIn = async () => {
    const result = await actions.signIn(email)
    if (result.ok) {
      setEmail("")
      setError(undefined)
      if (props.onEmailSent) {
        props.onEmailSent(email)
      }
    } else if (result.error) {
      setError(result.error)
    }
  }

  const providersLen = fider.settings.oauth.length

  if (!isCookieEnabled()) {
    return (
      <Message type="error">
        <h3 className="text-display">Cookies Required</h3>
        <p>Cookies are not enabled on your browser. Please enable cookies in your browser preferences to continue.</p>
      </Message>
    )
  }

  return (
    <div className="c-signin-control">
      {providersLen > 0 && (
        <>
          <div className="c-signin-control__oauth mb-2">
            {fider.settings.oauth.map((o) => (
              <React.Fragment key={o.provider}>
                <SocialSignInButton option={o} redirectTo={props.redirectTo} />
              </React.Fragment>
            ))}
          </div>
          {props.useEmail && <Divider />}
        </>
      )}

      {props.useEmail &&
        (showEmailForm ? (
          <div>
            <p>
              <Trans id="signin.message.email">Enter your email address to sign in</Trans>
            </p>
            <Form error={error}>
              <Input
                field="email"
                value={email}
                autoFocus={!device.isTouch()}
                onChange={setEmail}
                placeholder="yourname@example.com"
                suffix={
                  <Button type="submit" variant="primary" disabled={email === ""} onClick={signIn}>
                    <Trans id="action.signin">Sign in</Trans>
                  </Button>
                }
              />
            </Form>
            {!fider.session.tenant.isEmailAuthAllowed && (
              <p className="text-red-700 mt-1">
                <Trans id="signin.message.onlyadmins">Currently only allowed to sign in to an administrator account</Trans>
              </p>
            )}
          </div>
        ) : (
          <div>
            <p className="text-muted">
              <Trans id="signin.message.emaildisabled">
                Email authentication has been disabled by an administrator. If you have an administrator account and need to bypass this restriction, please{" "}
                <a href="#" className="text-bold" onClick={forceShowEmailForm}>
                  click here
                </a>
                .
              </Trans>
            </p>
          </div>
        ))}
    </div>
  )
}
