import * as Pages from "@fider/AsyncPages"

interface PageConfiguration {
  regex: RegExp
  component: any
  showHeader: boolean
}

export const route = (path: string, component: any, showHeader = true): PageConfiguration => {
  path = path.replace("/", "/").replace(":number", "\\d+").replace(":string", ".+").replace("*", "/?.*")

  const regex = new RegExp(`^${path}$`)
  return { regex, component, showHeader }
}

const defaultRoutes = [
  route("", Pages.AsyncHomePage),
  route("/posts/:number*", Pages.AsyncShowPostPage),
  route("/admin/members", Pages.AsyncManageMembersPage),
  route("/admin/tags", Pages.AsyncManageTagsPage),
  route("/admin/privacy", Pages.AsyncPrivacySettingsPage),
  route("/admin/export", Pages.AsyncExportPage),
  route("/admin/invitations", Pages.AsyncInvitationsPage),
  route("/admin/authentication", Pages.AsyncManageAuthenticationPage),
  route("/admin/advanced", Pages.AsyncAdvancedSettingsPage),
  route("/admin/webhooks", Pages.AsyncManageWebhooksPage),
  route("/admin", Pages.AsyncGeneralSettingsPage),
  route("/terms", Pages.AsyncLegalPage, false),
  route("/privacy", Pages.AsyncLegalPage, false),
  route("/signin", Pages.AsyncSignInPage, false),
  route("/signup", Pages.AsyncSignUpPage, false),
  route("/signin/verify", Pages.AsyncCompleteSignInProfilePage),
  route("/invite/verify", Pages.AsyncCompleteSignInProfilePage),
  route("/notifications", Pages.AsyncMyNotificationsPage),
  route("/settings", Pages.AsyncMySettingsPage),
  route("/oauth/:string/echo", Pages.AsyncOAuthEchoPage, false),
  route("/_design", Pages.AsyncDesignSystemPage),
]

export const resolveRootComponent = (path: string, routes: PageConfiguration[] = defaultRoutes): PageConfiguration => {
  if (path.length > 0 && path.charAt(path.length - 1) === "/") {
    path = path.substring(0, path.length - 1)
  }
  for (const entry of routes) {
    if (entry && entry.regex.test(path)) {
      return entry
    }
  }
  throw new Error(`Component not found for route ${path}.`)
}
