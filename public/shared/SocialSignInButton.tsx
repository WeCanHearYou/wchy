import * as React from 'react';

import { inject, injectables } from '../di';
import { Session } from '../services/Session';

interface SocialSignInButtonProps {
    provider: string;
    size: 'small' | 'normal';
}
interface SocialSignInButtonState {
    clicked: boolean;
}

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, SocialSignInButtonState> {

    @inject(injectables.Session)
    public session: Session;

    constructor() {
        super();
        this.state = {
            clicked: false
        };
    }

    public render() {
        const auth = this.session.get<any>('auth');
        const providerClassName = this.props.provider === 'google' ? 'google plus' : 'facebook';
        const providerDisplayName = this.props.provider === 'google' ? 'Google' : 'Facebook';
        const oauthUrl = `${auth.endpoint}/oauth/${this.props.provider}?redirect=${location.href}`;
        const cssClasses = `ui button
                            ${providerClassName}
                            ${this.state.clicked ? 'loading disabled' : ''}
                            ${this.props.size === 'small' ? 'circular icon' : 'fluid'}`;

        if (this.props.size === 'small') {
            return <a href={oauthUrl} className={cssClasses} onClick={() => this.setState({clicked: true})}>
                    <i className={providerClassName + ' icon'}></i>
                    </a>;
        } else {
            return  <a href={oauthUrl} className={cssClasses} onClick={() => this.setState({clicked: true})}>
                        <i className={providerClassName + ' icon'}></i>
                        Sign in with { providerDisplayName }
                    </a>;
        }
    }
}
