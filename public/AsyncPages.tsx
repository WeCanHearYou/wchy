import { lazy, ComponentType } from "react"

type LazyImport = () => Promise<{ default: ComponentType<any> }>

const MAX_RETRIES = 10
const INTERVAL = 500

const retry = (fn: LazyImport, retriesLeft = MAX_RETRIES): Promise<{ default: ComponentType<any> }> => {
  return new Promise((resolve, reject) => {
    fn()
      .then(resolve)
      .catch((err) => {
        setTimeout(() => {
          if (retriesLeft === 1) {
            reject(new Error(`${err} after ${MAX_RETRIES} retries`))
            return
          }
          retry(fn, retriesLeft - 1).then(resolve, reject)
        }, INTERVAL)
      })
  })
}

const load = (fn: LazyImport) => lazy(() => retry(() => fn()))

export const AsyncHomePage = load(
  () =>
    import(
      /* webpackChunkName: "Home.page" */
      "@fider/pages/Home/Home.page"
    )
)

export const AsyncShowPostPage = load(
  () =>
    import(
      /* webpackChunkName: "ShowPost.page" */
      "@fider/pages/ShowPost/ShowPost.page"
    )
)

export const AsyncManageMembersPage = load(
  () =>
    import(
      /* webpackChunkName: "ManageMembers.page" */
      "@fider/pages/Administration/pages/ManageMembers.page"
    )
)

export const AsyncManageTagsPage = load(
  () =>
    import(
      /* webpackChunkName: "ManageTags.page" */
      "@fider/pages/Administration/pages/ManageTags.page"
    )
)

export const AsyncPrivacySettingsPage = load(
  () =>
    import(
      /* webpackChunkName: "PrivacySettings.page" */
      "@fider/pages/Administration/pages/PrivacySettings.page"
    )
)

export const AsyncExportPage = load(
  () =>
    import(
      /* webpackChunkName: "Export.page" */
      "@fider/pages/Administration/pages/Export.page"
    )
)

export const AsyncInvitationsPage = load(
  () =>
    import(
      /* webpackChunkName: "Invitations.page" */
      "@fider/pages/Administration/pages/Invitations.page"
    )
)

export const AsyncManageAuthenticationPage = load(
  () =>
    import(
      /* webpackChunkName: "ManageAuthentication.page" */
      "@fider/pages/Administration/pages/ManageAuthentication.page"
    )
)

export const AsyncAdvancedSettingsPage = load(
  () =>
    import(
      /* webpackChunkName: "AdvancedSettings.page" */
      "@fider/pages/Administration/pages/AdvancedSettings.page"
    )
)

export const AsyncManageWebhooksPage = load(
  () =>
    import(
      /* webpackChunkName: "ManageWebhooks.page" */
      "@fider/pages/Administration/pages/ManageWebhooks.page"
    )
)

export const AsyncGeneralSettingsPage = load(
  () =>
    import(
      /* webpackChunkName: "GeneralSettings.page" */
      "@fider/pages/Administration/pages/GeneralSettings.page"
    )
)

export const AsyncSignInPage = load(
  () =>
    import(
      /* webpackChunkName: "SignIn.page" */
      "@fider/pages/SignIn/SignIn.page"
    )
)

export const AsyncSignUpPage = load(
  () =>
    import(
      /* webpackChunkName: "SignUp.page" */
      "@fider/pages/SignUp/SignUp.page"
    )
)

export const AsyncCompleteSignInProfilePage = load(
  () =>
    import(
      /* webpackChunkName: "CompleteSignInProfile.page" */
      "@fider/pages/CompleteSignInProfile/CompleteSignInProfile.page"
    )
)

export const AsyncMyNotificationsPage = load(
  () =>
    import(
      /* webpackChunkName: "MyNotifications.page" */
      "@fider/pages/MyNotifications/MyNotifications.page"
    )
)

export const AsyncMySettingsPage = load(
  () =>
    import(
      /* webpackChunkName: "MySettings.page" */
      "@fider/pages/MySettings/MySettings.page"
    )
)

export const AsyncOAuthEchoPage = load(
  () =>
    import(
      /* webpackChunkName: "OAuthEcho.page" */
      "@fider/pages/OAuthEcho/OAuthEcho.page"
    )
)

export const AsyncDesignSystemPage = load(
  () =>
    import(
      /* webpackChunkName: "DesignSystem.page" */
      "@fider/pages/DesignSystem/DesignSystem.page"
    )
)

export const AsyncLegalPage = load(
  () =>
    import(
      /* webpackChunkName: "Legal.page" */
      "@fider/pages/Legal/Legal.page"
    )
)
