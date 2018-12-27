import "./ImageUploader.scss";

import React from "react";
import { ValidationContext } from "./Form";
import { DisplayError, hasError } from "./DisplayError";
import { classSet, fileToBase64, uploadedImageURL } from "@fider/services";
import { Button, ButtonClickEvent, Modal } from "@fider/components";
import { FaRegImage } from "react-icons/fa";
import { ImageUpload } from "@fider/models";

interface ImageUploaderProps {
  instanceId?: string;
  field: string;
  label?: string;
  bkey?: string;
  disabled?: boolean;
  previewMaxWidth: number;
  onChange(state: ImageUpload, instanceId?: string): void;
}

interface ImageUploaderState extends ImageUpload {
  previewURL?: string;
  showModal: boolean;
}

export class ImageUploader extends React.Component<ImageUploaderProps, ImageUploaderState> {
  private fileSelector?: HTMLInputElement | null;

  constructor(props: ImageUploaderProps) {
    super(props);
    this.state = {
      upload: undefined,
      remove: false,
      showModal: false,
      previewURL: uploadedImageURL(this.props.bkey, this.props.previewMaxWidth)
    };
  }

  public fileChanged = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      const base64 = await fileToBase64(file);
      this.setState(
        {
          bkey: this.props.bkey,
          upload: {
            fileName: file.name,
            content: base64,
            contentType: file.type,
            action: "upload"
          },
          remove: false,
          previewURL: `data:${file.type};base64,${base64}`
        },
        () => {
          this.props.onChange(this.state, this.props.instanceId);
        }
      );
    }
  };

  public removeFile = async (e: ButtonClickEvent) => {
    if (this.fileSelector) {
      this.fileSelector.value = "";
    }

    this.setState(
      {
        bkey: this.props.bkey,
        remove: true,
        upload: undefined,
        previewURL: undefined
      },
      () => {
        this.props.onChange(
          {
            bkey: this.state.bkey,
            remove: this.state.remove,
            upload: this.state.upload
          },
          this.props.instanceId
        );
      }
    );
  };

  public selectFile = async (e: ButtonClickEvent) => {
    if (this.fileSelector) {
      this.fileSelector.click();
    }
  };

  private openModal = () => {
    this.setState({ showModal: true });
  };

  private closeModal = async () => {
    this.setState({ showModal: false });
  };

  private modal() {
    return (
      <Modal.Window isOpen={this.state.showModal} center={false} size="fluid">
        <Modal.Content>
          <img src={this.state.previewURL} />
        </Modal.Content>

        <Modal.Footer>
          <Button color="cancel" onClick={this.closeModal}>
            Close
          </Button>
        </Modal.Footer>
      </Modal.Window>
    );
  }

  public render() {
    const isUploading = !!this.state.upload;
    const hasFile = (!this.state.remove && this.props.bkey) || isUploading;

    const imgStyles: React.CSSProperties = {
      maxWidth: `${this.props.previewMaxWidth}px`
    };

    return (
      <ValidationContext.Consumer>
        {ctx => (
          <div
            className={classSet({
              "c-form-field": true,
              "c-image-upload": true,
              "m-error": hasError(this.props.field, ctx.error)
            })}
          >
            {this.modal()}
            <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>

            {hasFile && (
              <div className="preview">
                <img onClick={this.openModal} src={this.state.previewURL} style={imgStyles} />
                {!this.props.disabled && (
                  <Button onClick={this.removeFile} color="danger">
                    X
                  </Button>
                )}
              </div>
            )}

            <input ref={e => (this.fileSelector = e)} type="file" onChange={this.fileChanged} />
            <DisplayError fields={[this.props.field]} error={ctx.error} />
            {!hasFile && (
              <div className="c-form-field-wrapper">
                <Button onClick={this.selectFile} disabled={this.props.disabled}>
                  <FaRegImage />
                </Button>
              </div>
            )}
            {this.props.children}
          </div>
        )}
      </ValidationContext.Consumer>
    );
  }
}
