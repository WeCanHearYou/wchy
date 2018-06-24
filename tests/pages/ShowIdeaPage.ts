import { Browser, Page, WebComponent, TextInput, Button, findBy, elementIsVisible } from "../lib";

export class ShowIdeaPage extends Page {
  constructor(browser: Browser) {
    super(browser);
  }

  @findBy(".idea-header h1") public Title!: WebComponent;
  @findBy(".description") public Description!: WebComponent;
  @findBy(".c-segment.l-response .status-label") public Status!: WebComponent;
  @findBy(".c-segment.l-response .content") public ResponseText!: WebComponent;
  @findBy(".c-support-counter button") public SupportCounter!: WebComponent;
  @findBy(".c-comment-input #input-content") public CommentInput!: TextInput;
  @findBy(".c-comment-input .c-button.m-positive") public SubmitCommentButton!: Button;
  @findBy(".action-col .c-button.respond") public RespondButton!: Button;
  @findBy(".c-modal-window .c-response-form") public ResponseModal!: WebComponent;
  // @findMultipleBy(".c-comment-list .c-comment") public CommentList!: CommentList;
  // @findBy(".c-modal-window .c-response-form #input-status") private ResponseModalStatus!: DropDownList;
  @findBy(".c-modal-window .c-response-form #input-text") private ResponseModalText!: TextInput;
  @findBy(".c-modal-window .c-modal-footer .c-button.m-positive") private ResponseModalSubmitButton!: Button;

  @findBy(".action-col .c-button.edit") private Edit!: Button;
  @findBy("#input-title") private EditTitle!: TextInput;
  @findBy("#input-description") private EditDescription!: TextInput;
  @findBy(".action-col .c-button.save") private SaveEdit!: Button;
  @findBy(".action-col .c-button.cancel") private CancelEdit!: Button;

  public loadCondition() {
    return elementIsVisible(this.Title);
  }

  public async changeStatus(status: string, text: string): Promise<void> {
    // await this.ResponseModalStatus.selectByText(status);
    await this.ResponseModalText.clear();
    await this.ResponseModalText.type(text);
    await this.ResponseModalSubmitButton.click();
  }

  public async edit(newTitle: string, newDescription: string): Promise<void> {
    await this.Edit.click();
    await this.EditTitle.clear();
    await this.EditTitle.type(newTitle);
    await this.EditDescription.clear();
    await this.EditDescription.type(newDescription);
    await this.SaveEdit.click();
  }
}
