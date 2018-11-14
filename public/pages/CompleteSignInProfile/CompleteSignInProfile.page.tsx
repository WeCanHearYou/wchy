import React from "react";

import { HomePage, HomePageProps } from "../Home";
import { SignInPage } from "../SignIn";
import { Modal, Button, Form, Input, LegalFooter } from "@fider/components";
import { actions, Failure, querystring, Fider } from "@fider/services";

interface CompleteSignInProfilePageState {
  name: string;
  error?: Failure;
}

export class CompleteSignInProfilePage extends React.Component<HomePageProps, CompleteSignInProfilePageState> {
  private key: string;

  constructor(props: HomePageProps) {
    super(props);
    this.key = querystring.get("k");
    this.state = {
      name: ""
    };
  }

  private submit = async () => {
    const result = await actions.completeProfile(this.key, this.state.name);
    if (result.ok) {
      location.href = "/";
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  };

  private setName = (name: string) => {
    this.setState({ name });
  };

  public render() {
    return (
      <>
        <Modal.Window canClose={false} isOpen={true}>
          <Modal.Header>Complete your profile</Modal.Header>
          <Modal.Content>
            <p>Because this is your first sign in, please enter your name.</p>
            <Form error={this.state.error}>
              <Input
                field="name"
                onChange={this.setName}
                maxLength={100}
                onSubmit={this.submit}
                placeholder="Name"
                suffix={
                  <Button onClick={this.submit} color="positive" disabled={this.state.name === ""}>
                    Submit
                  </Button>
                }
              />
            </Form>
          </Modal.Content>
          <LegalFooter />
        </Modal.Window>
        {Fider.session.tenant.isPrivate
          ? React.createElement(SignInPage, this.props)
          : React.createElement(HomePage, this.props)}
      </>
    );
  }
}
