import { Browser, BrowserTab, mailgun, pageHasLoaded, ensure } from "../lib";
import { HomePage } from "../pages";

describe("E2E: Sign up with e-mail", () => {
  let browser: Browser;
  let tab: BrowserTab;

  beforeAll(async () => {
    browser = await Browser.launch();
    tab = await browser.newTab();
  });

  afterAll(async () => {
    await browser.close();
  });

  it("User can sign up using email", async () => {
    const now = new Date().getTime();

    // Action
    await tab.pages.signup.navigate();
    await tab.pages.signup.signInWithEmail(`Darth Vader ${now}`, `darthvader.fider@gmail.com`);
    await tab.pages.signup.signUpAs(`Selenium ${now}`, `selenium${now}`);

    const link = await mailgun.getLinkFromLastEmailTo(`selenium${now}`, `darthvader.fider@gmail.com`);

    await tab.pages.goTo(link);
    await tab.wait(pageHasLoaded(HomePage));

    // Assert
    await tab.pages.home.UserMenu.click();
    await ensure(tab.pages.home.UserName).textIs(`DARTH VADER ${now}`);
  });
});
