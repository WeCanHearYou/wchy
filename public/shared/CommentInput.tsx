import axios from 'axios';
import * as React from 'react';

import { Idea } from '../models';
import { DisplayError } from './Common';
import { SocialSignInList } from './SocialSignInList';

import { inject, injectables } from '../di';
import { Session } from '../services/Session';

interface CommentInputProps {
    idea: Idea;
}

interface CommentInputState {
    content: string;
    clicked: boolean;
    error?: Error;
}

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {

    @inject(injectables.Session)
    public session: Session;

    constructor() {
        super();
        this.state = {
          content: '',
          clicked: false
        };
    }

    public async submit() {
      this.setState({
        clicked: true,
        error: undefined
      });

      try {
        await axios.post(`/api/ideas/${this.props.idea.number}/comments`, {
          content: this.state.content
        });

        location.reload();
      } catch (ex) {
        this.setState({
          clicked: false,
          error: ex.response.data
        });
      }
    }

    public render() {
        const user = this.session.getCurrentUser();
        const buttonClasses = `ui blue labeled submit icon button ${this.state.clicked && 'loading disabled'}`;

        const addComment = user ? <form className="ui reply form">
          <DisplayError error={this.state.error} />
          <div className="field">
            <textarea onKeyUp={(e) => { this.setState({ content: e.currentTarget.value }); }}
                      placeholder="Leave your comment here..."></textarea>
          </div>
          <div className={ buttonClasses } onClick={async () => await this.submit()}>
            <i className="icon edit"></i> Add Comment
          </div>
        </form> :
        <div className="ui message">
          <div className="header">
            Please log in before leaving a comment.
          </div>
          <p>This only takes a second and you'll be good to go!</p>
          <SocialSignInList orientation="horizontal" size="small" />
        </div>;

        return addComment;
    }
}
