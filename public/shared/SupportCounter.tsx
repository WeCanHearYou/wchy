import * as React from 'react';
import { Idea, User, IdeaStatus } from '../models';
import { SocialSignInList } from './SocialSignInList';

import { inject, injectables } from '../di';
import { Session } from '../services/Session';
import { IdeaService } from '../services/IdeaService';

interface SupportCounterProps {
    user: User;
    idea: Idea;
}

interface SupportCounterState {
    supported: boolean;
    total: number;
}

export class SupportCounter extends React.Component<SupportCounterProps, SupportCounterState> {
    private elem: HTMLElement;
    private list: HTMLElement;

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.IdeaService)
    public service: IdeaService;

    constructor(props: SupportCounterProps) {
        super(props);
        const supportedIdeas = this.session.getArray<number>('supportedIdeas');

        this.state = {
          supported: props.user && supportedIdeas && supportedIdeas.indexOf(props.idea.id) >= 0,
          total: props.idea.totalSupporters
        };
    }

    public componentDidMount() {
        if (!this.props.user) {
            $(this.elem).popup({
                inline: true,
                hoverable: true,
                popup: this.list,
                on: 'click',
                position : 'bottom left'
            });
            return;
        }
    }

    public async supportOrUndo() {
        if (!this.props.user) {
            return;
        }

        const action = this.state.supported ? this.service.removeSupport : this.service.addSupport;

        const response = await action(this.props.idea.number);
        if (response.ok) {
            this.setState({
                supported: !this.state.supported,
                total: this.state.total + (this.state.supported ? -1 : 1)
            });
        } else {
            // TODO: handle this. we should have a global alert box
        }

    }

    public render() {
        const status = IdeaStatus.Get(this.props.idea.status);

        const support = <div className="ui mini violet inverted animated button"
                    onClick={async () => await this.supportOrUndo()}>
                    <div className="visible content">Want</div>
                    <div className="hidden content">
                        <i className="heart icon"></i>
                    </div>
                </div>;

        const undo = <div className="ui mini violet animated button"
                    onClick={async () => await this.supportOrUndo()}>
                    <div className="visible content"><i className="heart icon"></i></div>
                    <div className="hidden content">Undo</div>
                </div>;

        const disabled = <div className="ui disabled mini animated button">
                    <div className="visible content">~</div>
                </div>;

        return <div>
                    <div className="support-counter ui small statistics">
                        <div ref={(e) => { this.elem = e!; } } className="statistic">
                            <div className="value">
                                { this.state.total }
                            </div>
                            { status.closed ? disabled : this.state.supported ? undo : support }
                        </div>
                    </div>
                    <div ref={(e) => { this.list = e!; } } className="ui popup transition hidden">
                        <p className="header">Please log in to support this idea</p>
                        <SocialSignInList orientation="horizontal" size="small" />
                    </div>
                </div>;
    }
}
