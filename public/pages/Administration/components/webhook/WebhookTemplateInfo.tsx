import "./WebhookTemplateInfo.scss"

import { Button, Loader, Modal } from "@fider/components"
import { WebhookProperties } from "@fider/pages/Administration/components/webhook/WebhookProperties"
import React, { useEffect, useState } from "react"
import { WebhookType } from "@fider/models"
import { actions, StringObject } from "@fider/services"
import { VStack } from "@fider/components/layout"
import { HoverInfo } from "@fider/components/common/HoverInfo"

interface WebhookTemplateInfoProps {
  type: WebhookType
  isModalOpen: boolean
  onModalClose: () => void
}

interface FunctionSpecification {
  params: {
    type: string
    desc: string
    info?: string
    link?: string
  }[]
  description: string
  info?: string
  link?: string
}

const functions: StringObject<FunctionSpecification> = {
  stripHtml: {
    params: [{ type: "string", desc: "The input string containing HTML to strip" }],
    description: "Strip HTML tags from the input data",
    info: "It removes all tags to keep only content",
  },
  md5: {
    params: [{ type: "string", desc: "The input string to hash" }],
    description: "Hash text using md5 algorithm",
    info: "What is MD5?",
    link: "https://en.wikipedia.org/wiki/MD5",
  },
  lower: {
    params: [{ type: "string", desc: "The input string to lowercase" }],
    description: "Lowercase text",
  },
  upper: {
    params: [{ type: "string", desc: "The input string to uppercase" }],
    description: "Uppercase text",
  },
  markdown: {
    params: [{ type: "string", desc: "The input string containing Markdown to parse" }],
    description: "Parse Markdown to HTML from the input data",
    info: "When parsing, input is sanitized from HTML to prevent XSS attacks",
  },
  format: {
    params: [
      {
        type: "string",
        desc: "The date format according to Go specifications",
        info: "See Go time format",
        link: "https://yourbasic.org/golang/format-parse-string-time-date-example/#standard-time-and-date-formats",
      },
      { type: "time", desc: "The time to format" },
    ],
    description: "Format the given date and time",
  },
  quote: {
    params: [{ type: "string", desc: "The input string to quote" }],
    description: "Enquote a string and escape inner special characters",
    info: "You should use this function when using a value as a JSON field",
  },
  urlquery: {
    params: [{ type: "string", desc: "The input string to encode as URL query value" }],
    description: "Encode a string into a valid URL query element",
    info: "You should use this function when using a value in URL",
  },
}
const textExample = 'A new post entitled "{{ .post_title }}" has been created by {{ .author_name }}.'
const jsonExample = `{
  "title": {{ quote .post_title }},
  "content": {{ quote .post_description }},
  "user": {{ quote .author_name }}
}`

export const WebhookTemplateInfo = (props: WebhookTemplateInfoProps) => {
  const [properties, setProperties] = useState<StringObject | null>()

  useEffect(() => {
    let mounted = true
    setProperties(undefined)
    actions.getWebhookHelp(props.type).then((result) => mounted && setProperties(result.ok ? result.data : null))
    return () => {
      mounted = false
    }
  }, [props.type])

  return (
    <Modal.Window className="c-webhook-templateinfo" isOpen={props.isModalOpen} onClose={props.onModalClose} size="large">
      <Modal.Header>Webhook template formatting help</Modal.Header>
      <Modal.Content>
        <VStack spacing={4} divide>
          <div>
            <h3 className="text-title mb-1">What is a template?</h3>
            <p>
              The template engine used is the native Go <code>text/template</code> package. The simpliest way to create a template is to write your text, and
              insert the property name, prefixed by a dot, enclosed in double braces with spaces within (wierd? complete!). Example:
            </p>
            <pre className="text-left">{textExample}</pre>
            <p>
              When using a structured text format such as JSON, you may need to wrap those properties into quotes. But be careful: you must escape values by
              yourself by calling the appropriate function depending of your situation. Example:
            </p>
            <pre className="text-left">{jsonExample}</pre>
            <Button href="https://pkg.go.dev/text/template" target="_blank" variant="primary">
              Open official Go documentation about templates
            </Button>
          </div>
          {properties === undefined ? (
            <p className="text-muted">Failed to load help data</p>
          ) : properties === null ? (
            <Loader text="Loading help data" />
          ) : (
            <>
              <div>
                <h3 className="text-title mb-1">Properties</h3>
                <WebhookProperties properties={properties} propsName="Property name" valueName="Example value" />
              </div>
              <div>
                <h3 className="text-title mb-1">Functions</h3>
                <table className="c-webhook-properties">
                  <thead>
                    <tr>
                      <th rowSpan={2}>Function</th>
                      <th rowSpan={2}>Description</th>
                      <th colSpan={2}>Parameters</th>
                    </tr>
                    <tr>
                      <th className="text-xs">Type</th>
                      <th className="text-xs">Description</th>
                    </tr>
                  </thead>
                  <tbody>
                    {Object.entries(functions).map(([func, spec]) =>
                      spec.params.map((param, j) => (
                        <tr key={func}>
                          {j === 0 && (
                            <>
                              <td rowSpan={spec.params.length} className="c-webhook-properties__prop">
                                {func}
                              </td>
                              <td rowSpan={spec.params.length} className="c-webhook-properties__val">
                                {spec.description}
                                {spec.info && <HoverInfo text={spec.info} href={spec.link} target="_blank" />}
                              </td>
                            </>
                          )}
                          <td className="c-webhook-properties__prop">{param.type}</td>
                          <td className="c-webhook-properties__val">
                            {param.desc}
                            {param.info && <HoverInfo text={param.info} href={param.link} target="_blank" />}
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </>
          )}
        </VStack>
      </Modal.Content>
      <Modal.Footer>
        <Button variant="tertiary" onClick={props.onModalClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal.Window>
  )
}
