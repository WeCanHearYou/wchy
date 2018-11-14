import React from "react";
import ReactDOM from "react-dom";
import { resolveRootComponent } from "@fider/router";
import { Header, Footer } from "@fider/components/common";
import { ErrorBoundary } from "@fider/components";
import { classSet, Fider, actions, navigator } from "@fider/services";
import { IconContext } from "react-icons";

import "@fider/assets/styles/main.scss";

const logProductionError = (err: Error) => {
  if (Fider.isProduction()) {
    console.error(err); // tslint:disable-line
    actions.logError(`react.ErrorBoundary: ${err.message}`, err);
  }
};

window.addEventListener("unhandledrejection", (evt: PromiseRejectionEvent) => {
  if (evt.reason instanceof Error) {
    actions.logError(`window.unhandledrejection: ${evt.reason.message}`, evt.reason);
  } else if (evt.reason) {
    actions.logError(`window.unhandledrejection: ${evt.reason.toString()}`);
  }
});

window.addEventListener("error", (evt: ErrorEvent) => {
  if (evt.error && evt.colno > 0 && evt.lineno > 0) {
    actions.logError(`window.error: ${evt.message}`, evt.error);
  }
});

(() => {
  let fider;

  if (!navigator.isBrowserSupported()) {
    navigator.goTo("/browser-not-supported");
    return;
  }

  fider = Fider.initialize();

  __webpack_nonce__ = fider.session.contextID;
  __webpack_public_path__ = `${fider.settings.globalAssetsURL}/assets/`;

  const config = resolveRootComponent(location.pathname);
  document.body.className = classSet({
    "is-authenticated": fider.session.isAuthenticated,
    "is-staff": fider.session.isAuthenticated && fider.session.user.isCollaborator
  });
  ReactDOM.render(
    <ErrorBoundary onError={logProductionError}>
      <IconContext.Provider value={{ className: "icon" }}>
        {config.showHeader && <Header />}
        {React.createElement(config.component, fider.session.props)}
        {config.showHeader && <Footer />}
      </IconContext.Provider>
    </ErrorBoundary>,
    document.getElementById("root")
  );
})();
