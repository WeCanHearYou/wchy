import * as React from "react";

import {
  HomePage,
  SignInPage,
  SignUpPage,
  ManageMembersPage,
  CompleteSignInProfilePage,
  PrivacySettingsPage,
  InvitationsPage,
  GeneralSettingsPage,
  ManageTagsPage,
  ShowIdeaPage,
  MySettingsPage,
  MyNotificationsPage
} from "@fider/pages";

interface PageConfiguration {
  id: string;
  regex: RegExp;
  component: any;
  showHeader: boolean;
}

const route = (path: string, component: any, id: string, showHeader: boolean): PageConfiguration => {
  path = path
    .replace("/", "/")
    .replace(":number", "\\d+")
    .replace("*", "/?.*");

  const regex = new RegExp(`^${path}$`);
  return { regex, component, id, showHeader };
};

const pathRegex = [
  route("", HomePage, "p-home", true),
  route("/ideas/:number*", ShowIdeaPage, "p-show-idea", true),
  route("/admin/members", ManageMembersPage, "p-admin-members", true),
  route("/admin/tags", ManageTagsPage, "p-admin-tags", true),
  route("/admin/privacy", PrivacySettingsPage, "p-admin-privacy", true),
  route("/admin/invitations", InvitationsPage, "p-admin-invitations", true),
  route("/admin", GeneralSettingsPage, "p-admin-general", true),
  route("/signin", SignInPage, "p-signin", false),
  route("/signup", SignUpPage, "p-signup", false),
  route("/signin/verify", CompleteSignInProfilePage, "p-home", true),
  route("/invite/verify", CompleteSignInProfilePage, "p-signin", false),
  route("/notifications", MyNotificationsPage, "p-my-notifications", true),
  route("/settings", MySettingsPage, "p-my-settings", true)
];

export const resolveRootComponent = (path: string): PageConfiguration => {
  if (path.length > 0 && path.charAt(path.length - 1) === "/") {
    path = path.substring(0, path.length - 1);
  }
  for (const entry of pathRegex) {
    if (entry.regex.test(path)) {
      return entry;
    }
  }
  throw new Error(`Component not found for route ${path}.`);
};
