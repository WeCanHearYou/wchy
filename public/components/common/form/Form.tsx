import "./Form.scss";

import React from "react";
import { Failure, classSet } from "@fider/services";
import { DisplayError } from "@fider/components";

interface ValidationContext {
  error?: Failure;
}

interface FormProps {
  className?: string;
  size?: "mini" | "normal";
  error?: Failure;
  onSubmit?: (ev: React.FormEvent<HTMLFormElement>) => void;
}

export const ValidationContext = React.createContext<ValidationContext>({});

export const Form: React.StatelessComponent<FormProps> = props => {
  const className = classSet({
    "c-form": true,
    [props.className!]: props.className,
    [`m-${props.size}`]: props.size
  });

  return (
    <form autoComplete="off" className={className} onSubmit={props.onSubmit}>
      <DisplayError error={props.error} />
      <ValidationContext.Provider value={{ error: props.error }}>{props.children}</ValidationContext.Provider>
    </form>
  );
};
