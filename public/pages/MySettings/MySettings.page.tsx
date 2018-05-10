import "./MySettings.page.scss";

import * as React from "react";

import { Modal, Form, DisplayError, Button, Gravatar } from "@fider/components/common";
import { NotificationSettings } from "./";

import { CurrentUser, UserSettings } from "@fider/models";
import { Failure, actions } from "@fider/services";

interface MySettingsPageState {
  showModal: boolean;
  name: string;
  newEmail: string;
  changingEmail: boolean;
  error?: Failure;
  settings: UserSettings;
}

interface MySettingsPageProps {
  user: CurrentUser;
  settings: UserSettings;
}

export class MySettingsPage extends React.Component<MySettingsPageProps, MySettingsPageState> {
  constructor(props: MySettingsPageProps) {
    super(props);
    this.state = {
      showModal: false,
      changingEmail: false,
      newEmail: "",
      name: this.props.user.name,
      settings: this.props.settings
    };
  }

  private async confirm() {
    const result = await actions.updateUserSettings(this.state.name, this.state.settings);
    if (result.ok) {
      location.reload();
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  }

  private async submitNewEmail() {
    const result = await actions.changeUserEmail(this.state.newEmail);
    if (result.ok) {
      this.setState({
        error: undefined,
        changingEmail: false,
        showModal: true
      });
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  }

  public render() {
    return (
      <div id="p-my-settings" className="page ui container">
        <Modal.Window isOpen={this.state.showModal} canClose={true} center={true}>
          <Modal.Header>Confirm your new email</Modal.Header>
          <Modal.Content>
            <div>
              <p>
                We have just sent a confirmation link to <b>{this.state.newEmail}</b>. <br /> Click the link to update
                your email.
              </p>
              <p>
                <a href="#" onClick={() => this.setState({ showModal: false })}>
                  OK
                </a>
              </p>
            </div>
          </Modal.Content>
        </Modal.Window>

        <h2 className="ui header">
          <i className="circular id badge icon" />
          <div className="content">
            Settings
            <div className="sub header">Manage your profile settings</div>
          </div>
        </h2>

        <div className="ui grid">
          <div className="ten wide computer sixteen wide mobile column">
            <div className="ui form">
              <div className="field">
                <label htmlFor="email">Avatar</label>
                <p>
                  <Gravatar user={this.props.user} />
                </p>
                <div className="info">
                  <p>
                    This site uses{" "}
                    <a href="https://en.gravatar.com/" target="blank">
                      Gravatar
                    </a>{" "}
                    to display profile avatars. <br />
                    A letter avatar based on your name is generated for profiles without a Gravatar.
                  </p>
                </div>
              </div>
              <DisplayError fields={["email"]} error={this.state.error} />
              <div className="field">
                <label htmlFor="email">
                  Email <span className="info">Your email is private and will never be displayed to anyone.</span>
                </label>
                {this.state.changingEmail ? (
                  <>
                    <p>
                      <input
                        id="new-email"
                        type="text"
                        style={{ maxWidth: "200px", marginRight: "10px" }}
                        maxLength={200}
                        placeholder={this.props.user.email}
                        value={this.state.newEmail}
                        onChange={e => this.setState({ newEmail: e.currentTarget.value })}
                      />
                      <Button color="positive" size="mini" onClick={async () => await this.submitNewEmail()}>
                        Confirm
                      </Button>
                      <Button
                        size="mini"
                        onClick={async () =>
                          this.setState({
                            changingEmail: false,
                            error: undefined
                          })
                        }
                      >
                        Cancel
                      </Button>
                    </p>
                  </>
                ) : (
                  <p>
                    {this.props.user.email ? (
                      <b>{this.props.user.email}</b>
                    ) : (
                      <span className="info">Your account doesn't have an email.</span>
                    )}
                    <span className="ui info clickable" onClick={() => this.setState({ changingEmail: true })}>
                      change
                    </span>
                  </p>
                )}
              </div>

              <DisplayError fields={["name"]} error={this.state.error} />
              <div className="field">
                <label htmlFor="name">Name</label>
                <input
                  id="name"
                  type="text"
                  maxLength={100}
                  value={this.state.name}
                  onChange={e => this.setState({ name: e.currentTarget.value })}
                />
              </div>

              <NotificationSettings
                user={this.props.user}
                settings={this.props.settings}
                settingsChanged={settings => this.setState({ settings })}
              />

              <div className="field">
                <Button color="positive" onClick={async () => await this.confirm()}>
                  Confirm
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
